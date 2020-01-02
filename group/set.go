package group

type ElementsSet struct {
	data map[interface{}]int64
}

func NewElementsSet() *ElementsSet {
	return &ElementsSet{
		data: make(map[interface{}]int64),
	}
}

func (es *ElementsSet) Put(key GroupElement, value int64) {
	es.data[key.Hash()] = value
}

func (es *ElementsSet) Get(key GroupElement) (value int64, ok bool) {
	value, ok = es.data[key.Hash()]
	return value, ok
}

func (es *ElementsSet) Exists(key GroupElement) bool {
	_, ok := es.data[key.Hash()]
	return ok
}
