package shanks

import (
	"github.com/pkg/errors"
	"math"
	"shanks-algorithm/group"
)

type ShanksAlgorithm struct {
	a, b 	group.GroupElement
	g		group.Group
	es		*group.ElementsSet
	done	bool
	x    	int64
}

func NewShanksAlgorithm(a, b group.GroupElement, g group.Group) *ShanksAlgorithm {
	return &ShanksAlgorithm{
		a: a,
		b: b,
		g: g,
		es: group.NewElementsSet(),
		done: false,
	}
}

func (sa *ShanksAlgorithm) Execute() (x int64, err error) {
	if sa.done {
		return sa.x, nil
	}

	if sa.b.IsNeutralElement() {
		sa.done = true
		sa.x = 0
		return 0, nil
	}

	if sa.a.IsNeutralElement() && !sa.b.IsNeutralElement() {
		return 0, errors.New("wrong input: a is neutral, but b is not")
	}

	l := sa.g.GroupOrder()

	m := int64(math.Ceil(math.Sqrt(float64(l))))
	c := sa.a.Pow(m)
	for i := int64(1); i <= m; i++ {
		r := c.Pow(i)
		//fmt.Println(r)
		sa.es.Put(r, i)
	}
	//fmt.Println()

	found := false
	for j := int64(1); j <= m; j++ {
		t := sa.b.Add(sa.a.Pow(j))
		//fmt.Println(t)
		if i, ok := sa.es.Get(t); ok {
			x = (m * i - j) % l
			if sa.a.Pow(x).Equal(sa.b) {
				found = true
				break
			}
		}
	}

	if !found {
		return 0, errors.New("no matches found")
	}

	sa.done = true
	sa.x = x

	return x, nil
}


