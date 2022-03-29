package keys

import (
	"fmt"
	"github.com/99designs/keyring"
	"io"
	"path/filepath"
)

const (
	keyringTestDirName = "keyring-test"
)

func NewTestKeyring(
	appName, backend, rootDir string, userInput io.Reader, opts ...KeybaseOption,
) (Keybase, error) {

	var db keyring.Keyring
	var err error
	var config keyring.Config

	switch backend {
	case BackendTest:
		config = newTestBackendKeyringConfig(appName, rootDir)
	case BackendFile:
		config = newFileBackendKeyringConfig(appName, rootDir, userInput)
	case BackendOS:
		config = lkbToKeyringConfig(appName, rootDir, userInput, false)
	case BackendKWallet:
		config = newKWalletBackendKeyringConfig(appName, rootDir, userInput)
	case BackendPass:
		config = newPassBackendKeyringConfig(appName, rootDir, userInput)
	case BackendMemory:
		return NewInMemory(opts...), err
	default:
		return nil, fmt.Errorf("unknown keyring backend %v", backend)
	}
	db, err = keyring.Open(config)
	if err != nil {
		return nil, err
	}

	return newKeyringKeybase(db, config.FileDir, opts...), nil
}

func newTestBackendKeyringConfig(appName, dir string) keyring.Config {
	return keyring.Config{
		AllowedBackends: []keyring.BackendType{keyring.FileBackend},
		ServiceName:     appName,
		FileDir:         filepath.Join(dir, keyringTestDirName),
		FilePasswordFunc: func(_ string) (string, error) {
			return "test", nil
		},
	}
}
