package ethkeystore

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

type encryptedKeyJSONV3ForWeb3 struct {
	Address string              `json:"address"`
	Crypto  keystore.CryptoJSON `json:"crypto"`
	Id      string              `json:"id"`
	Version int                 `json:"version"`
}

// DecryptKey decrypts a key from a json blob, returning the private key itself.
func DecryptKeyForWeb3(keyjson []byte, auth string) (*keystore.Key, error) {
	// Parse the json into a simple map to fetch the key version
	m := make(map[string]interface{})
	if err := json.Unmarshal(keyjson, &m); err != nil {
		return nil, err
	}
	// Depending on the version try to parse one way or another
	var (
		keyBytes, keyId []byte
		err             error
	)
	if m["version"] == float64(3) {
		k := new(encryptedKeyJSONV3ForWeb3)
		if err := json.Unmarshal(keyjson, k); err != nil {
			return nil, err
		}
		keyBytes, keyId, err = decryptKeyV3ForWeb3(k, auth)
	} else {
		return nil, fmt.Errorf("the veison must be equal to 3, got: %s", m["version"])
	}
	// Handle any decryption errors and return the key
	if err != nil {
		return nil, err
	}
	key := ethcrypto.ToECDSAUnsafe(keyBytes)
	id, err := uuid.FromBytes(keyId)
	if err != nil {
		return nil, err
	}
	return &keystore.Key{
		Id:         id,
		Address:    ethcrypto.PubkeyToAddress(key.PublicKey),
		PrivateKey: key,
	}, nil
}

func decryptKeyV3ForWeb3(keyProtected *encryptedKeyJSONV3ForWeb3, auth string) (keyBytes []byte, keyId []byte, err error) {
	if keyProtected.Version != 3 {
		return nil, nil, fmt.Errorf("version not supported: %v", keyProtected.Version)
	}
	keyUUID, err := uuid.Parse(keyProtected.Id)
	if err != nil {
		return nil, nil, err
	}
	keyId = keyUUID[:]
	plainText, err := DecryptDataV3ForWeb3(keyProtected.Crypto, auth)
	if err != nil {
		return nil, nil, err
	}
	return plainText, keyId, err
}

func DecryptDataV3ForWeb3(cryptoJson keystore.CryptoJSON, auth string) ([]byte, error) {
	if cryptoJson.Cipher != "aes-128-ctr" {
		return nil, fmt.Errorf("cipher not supported: %v", cryptoJson.Cipher)
	}
	mac, err := hex.DecodeString(cryptoJson.MAC)
	if err != nil {
		return nil, err
	}

	iv, err := hex.DecodeString(cryptoJson.CipherParams.IV)
	if err != nil {
		return nil, err
	}

	cipherText, err := hex.DecodeString(cryptoJson.CipherText)
	if err != nil {
		return nil, err
	}

	derivedKey, err := getKDFKey(cryptoJson, auth)
	if err != nil {
		return nil, err
	}

	bufferValue := bytes.Buffer{}
	bufferValue.Write(derivedKey[16:32])
	bufferValue.Write(cipherText)
	calculatedMAC := sha256.Sum256(bufferValue.Bytes())
	if !bytes.Equal(calculatedMAC[:], mac) {
		return nil, keystore.ErrDecrypt
	}

	plainText, err := aesCTRXOR(derivedKey[:16], cipherText, iv)
	if err != nil {
		return nil, err
	}
	return plainText, err
}

func getKDFKey(cryptoJSON keystore.CryptoJSON, auth string) ([]byte, error) {
	authArray := []byte(auth)
	salt, err := hex.DecodeString(cryptoJSON.KDFParams["salt"].(string))
	if err != nil {
		return nil, err
	}
	dkLen := ensureInt(cryptoJSON.KDFParams["dklen"])

	if cryptoJSON.KDF == "scrypt" {
		n := ensureInt(cryptoJSON.KDFParams["n"])
		r := ensureInt(cryptoJSON.KDFParams["r"])
		p := ensureInt(cryptoJSON.KDFParams["p"])
		return scrypt.Key(authArray, salt, n, r, p, dkLen)
	} else if cryptoJSON.KDF == "pbkdf2" {
		c := ensureInt(cryptoJSON.KDFParams["c"])
		prf := cryptoJSON.KDFParams["prf"].(string)
		if prf != "hmac-sha256" {
			return nil, fmt.Errorf("unsupported PBKDF2 PRF: %s", prf)
		}
		key := pbkdf2.Key(authArray, salt, c, dkLen, sha256.New)
		return key, nil
	}

	return nil, fmt.Errorf("unsupported KDF: %s", cryptoJSON.KDF)
}

// TODO: can we do without this when unmarshalling dynamic JSON?
// why do integers in KDF params end up as float64 and not int after
// unmarshal?
func ensureInt(x interface{}) int {
	res, ok := x.(int)
	if !ok {
		res = int(x.(float64))
	}
	return res
}

func aesCTRXOR(key, inText, iv []byte) ([]byte, error) {
	// AES-128 is selected due to size of encryptKey.
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	outText := make([]byte, len(inText))
	stream.XORKeyStream(outText, inText)
	return outText, err
}
