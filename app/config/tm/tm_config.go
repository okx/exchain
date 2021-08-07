package tm

type TmConfig interface {
	GetTpb() uint64
}


var Config TmConfig


func SetConfig(c TmConfig) {
	Config = c
}