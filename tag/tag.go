package tag

type Tag struct {
	GroupNumber   uint16
	ElementNumber uint16
}

func (t1 Tag) Equal(t2 Tag) bool {
	return t1.GroupNumber == t2.GroupNumber && t1.ElementNumber == t2.ElementNumber
}
