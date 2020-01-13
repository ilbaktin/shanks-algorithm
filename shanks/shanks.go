package shanks

import (
	"github.com/pkg/errors"
	"math"
	"shanks-algorithm/group"
	"sync"
	"sync/atomic"
)

type ShanksAlgorithm struct {
	a, b group.GroupElement
	g    group.Group
	es   *group.ElementSetMap
	done bool
	x    int64
}

func NewShanksAlgorithm(a, b group.GroupElement, g group.Group) *ShanksAlgorithm {
	return &ShanksAlgorithm{
		a:    a,
		b:    b,
		g:    g,
		es:   group.NewElementsSetMap(),
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
			x = (m*i - j) % l
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

func (sa *ShanksAlgorithm) ExecuteParallel(threadsNum int) (x int64, err error) {
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
	ranges := splitRange(m, int64(threadsNum))
	elemSets := make([]*group.ElementSetMap, len(ranges))
	wg := sync.WaitGroup{}
	wg.Add(len(ranges))
	for idx, rng := range ranges {
		go func(idx int, rng int64Range) {
			elemSet := group.NewElementsSetMap()
			elemSets[idx] = elemSet
			for i := rng.start + 1; i <= rng.end; i++ {
				r := c.Pow(i)
				//fmt.Println(r)
				elemSet.Put(r, i)
			}
			wg.Done()
		}(idx, rng)
	}
	wg.Wait()

	for _, es := range elemSets {
		sa.es.MergeFrom(es)
	}

	var found int32 = 0 // 0 if solution not found, 1 if found
	//var checkPeriod int64 = 10

	wg = sync.WaitGroup{}
	wg.Add(len(ranges))
	solution := int64(0)
	for _, rng := range ranges {

		go func(rng int64Range) {
			for j := rng.start + 1; j <= rng.end; j++ {
				if found == 1 {
					break
				}
				t := sa.b.Add(sa.a.Pow(j))
				if i, ok := sa.es.Get(t); ok {
					x := (m*i - j) % l
					if sa.a.Pow(x).Equal(sa.b) {
						//fmt.Printf("Found solution x=%d\n", x)
						solution = x
						atomic.AddInt32(&found, 1)
						break
					}
				}
			}
			wg.Done()
		}(rng)
	}

	wg.Wait()
	//fmt.Printf("%#v\n", sa.es)

	if found > 0 {
		return solution, nil
	} else {
		return 0, errors.New("no matches found")
	}
}
