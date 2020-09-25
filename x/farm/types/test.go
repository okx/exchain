package types

// TODO: remove it later
type TestStruct struct {
	Reminder string
}

func (tp TestStruct) String() string {
	return tp.Reminder
}

func NewTestStruct(s string) TestStruct {
	return TestStruct{s}
}
