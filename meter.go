package metrics

type Meter interface {
	Count() int64
	Mark(int64)
	Rate1()
	Rate5()
	Rate15()
	RateMean()
}
