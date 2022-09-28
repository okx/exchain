package types

func (b Block) To() CM40Block {
	ret := CM40Block{}
	return ret.From(b)
}
func (c Commit) To() IBCCommit {
	ret := IBCCommit{}
	return ret.From(c)
}
