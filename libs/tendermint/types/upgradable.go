package types

type UpgradeAble interface {
	Upgrade() interface{}
}
