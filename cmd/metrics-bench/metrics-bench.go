package main

import (
	"fmt"
	"metrics"
	"time"
)

func main() {
	r := metrics.NewRegistry()
	for i := 0; i < 250; i++ {
		r.RegisterMeter(fmt.Sprintf("%d", i), metrics.NewMeter())
	}
	time.Sleep(600e9)
}
