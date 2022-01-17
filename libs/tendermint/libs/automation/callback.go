package automation

var (
	callBackroleAction map[string]*callBackAction
)

type callBackAction struct {
	f func(data ...interface{})
}

func CallBack(height int64, round int) {
	if !enableRoleTest {
		return
	}
	if act, ok := callBackroleAction[actionKey(height, round)]; ok {
		act.f()
	}
}
func PrerunCallBackWithTimeOut(height int64, round int) {
	PrerunTimeOut(height, round)
	CallBack(height, round)
}

func RegisterActionCallBack(height int64, round int, f func(data ...interface{})) {
	key := actionKey(height, round)
	act, ok := callBackroleAction[key]
	if ok {
		panic("duplicate action")
	}
	act = &callBackAction{
		f: f,
	}
	callBackroleAction[key] = act
}
