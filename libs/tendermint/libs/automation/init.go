package automation

func init() {
	once.Do(func() {
		roleAction = make(map[string]*action)
		callBackroleAction = make(map[string]*callBackAction)
	})
}
