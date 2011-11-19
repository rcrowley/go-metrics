package metrics

type Gauge interface {
	Update(int64)
	Value() int64
}

type gauge struct {
	in, out chan int64
}

func NewGauge() Gauge {
	g := &gauge{make(chan int64), make(chan int64)}
	go g.arbiter()
	return g
}

func (g *gauge) Update(v int64) {
	g.in <- v
}

func (g *gauge) Value() int64 {
	return <-g.out
}

func (g *gauge) arbiter() {
	var value int64
	for {
		select {
		case v := <-g.in: value = v
		case g.out <- value:
		}
	}
}
