package group

type ElementSet interface {
	Put(key GroupElement, value int64)
	Get(key GroupElement) (value int64, ok bool)
	Exists(key GroupElement) bool
}

type ElementSetMap struct {
	data map[interface{}]int64
}

func NewElementsSetMap() *ElementSetMap {
	return &ElementSetMap{
		data: make(map[interface{}]int64),
	}
}

func (es *ElementSetMap) Put(key GroupElement, value int64) {
	es.data[key.Hash()] = value
}

func (es *ElementSetMap) Get(key GroupElement) (value int64, ok bool) {
	hash := key.Hash()
	value, ok = es.data[hash]
	return value, ok
}

func (es *ElementSetMap) Exists(key GroupElement) bool {
	_, ok := es.data[key.Hash()]
	return ok
}

func (es *ElementSetMap) MergeFrom(other *ElementSetMap) {
	for k, v := range other.data {
		es.data[k] = v
	}
}

type ElementSetSyncMap struct {
	data map[interface{}]int64
}

func NewElementsSetSyncMap() *ElementSetSyncMap {
	return &ElementSetSyncMap{
		data: make(map[interface{}]int64),
	}
}

func (es *ElementSetSyncMap) Put(key GroupElement, value int64) {
	es.data[key.Hash()] = value
}

func (es *ElementSetSyncMap) Get(key GroupElement) (value int64, ok bool) {
	value, ok = es.data[key.Hash()]
	return value, ok
}

func (es *ElementSetSyncMap) Exists(key GroupElement) bool {
	_, ok := es.data[key.Hash()]
	return ok
}
