package metrics

type Counter interface {
	Clear()
	Count() int64
	Dec(int64)
	Inc(int64)
}

type counter struct {
	in, out chan int64
	reset chan bool
}

func NewCounter() Counter {
	c := &counter{make(chan int64), make(chan int64), make(chan bool)}
	go c.arbiter()
	return c
}

func (c *counter) Clear() {
	c.reset <- true
}

func (c *counter) Count() int64 {
	return <-c.out
}

func (c *counter) Dec(i int64) {
	c.in <- -i
}

func (c *counter) Inc(i int64) {
	c.in <- i
}

func (c *counter) arbiter() {
	var count int64
	for {
		select {
		case i := <-c.in: count += i
		case c.out <- count:
		case <-c.reset: count = 0
		}
	}
}
