package seven

import (
	"context"
	"fmt"
	"time"
)

func ExampleSeven() {
	or := orClose
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("fone after %v", time.Since(start))
}

func orClose(channels ...<-chan interface{}) <-chan interface{} {
	done := make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, ch := range channels {
		go func(ctx context.Context, channel <-chan interface{}, done chan bool) {
			for {
				_, ok := <-channel
				// если закрыт, то завершаем
				if !ok {
					done <- true
				}
			}
		}(ctx, ch, done)
	}
	<-done
	cancel()
	res := make(chan interface{})
	close(res)
	return res
}
