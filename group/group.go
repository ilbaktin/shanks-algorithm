package group

type Group interface {
	GroupOrder() int64
	//ElementOrder(element GroupElement) (int64, error)
}

type GroupElement interface {
	Add(GroupElement) GroupElement
	Sub(GroupElement) GroupElement
	Pow(int64) GroupElement
	Equal(GroupElement) bool
	IsNeutralElement() bool
	Hash() interface{}
}
