package types

type WrapCMTx struct {
	Tx    Tx     `json:"tx" yaml:"tx"`
	Nonce uint64 `json:"nonce" yaml:"nonce"`
}

func (wtx *WrapCMTx) GetTx() Tx {
	if wtx != nil {
		return wtx.Tx
	}
	return nil
}

func (wtx *WrapCMTx) GetNonce() uint64 {
	if wtx != nil {
		return wtx.Nonce
	}
	return 0
}
