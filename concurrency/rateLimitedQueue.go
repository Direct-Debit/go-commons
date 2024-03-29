package concurrency

import (
	"github.com/Direct-Debit/go-commons/stdext"
	"time"
)

type RateLimitedQueueFull struct{}

func (r *RateLimitedQueueFull) Error() string {
	return "the rate limited queue is full"
}

// RateLimitedQueue can reduce the rate of incoming requests, and handle them in a timed and orderly fashion.
// This is useful for:
//   - Handling bursts of traffic.
//   - Add predictability to unpredictable traffic patterns.
//   - Making sure that you don't overload third party services.
type RateLimitedQueue[T any] struct {
	inChannel  chan T
	outChannel chan T
}

// RateLimitedQueueConfig allows you to specify configuration when creating a new rate limiter object.
//   - BufferSize: Allows you to specify a buffer size for the rate limiter. If the buffer is full adding jobs to the limiter will block execution.
//   - MaxConcurrent: Specify the concurrency limit of the resulting rate. Minimum value is 1, which indicates no concurrency.
//   - Rate: The resulting rate. Minimum value is one NanoSecond.
type RateLimitedQueueConfig struct {
	BufferSize    int
	MaxConcurrent int
	Rate          time.Duration
}

// NewRateLimitedQueue generates and starts up a new rate limiter.
func NewRateLimitedQueue[T any](cfg RateLimitedQueueConfig) *RateLimitedQueue[T] {
	cfg.BufferSize = stdext.Max(cfg.BufferSize, 0)
	cfg.MaxConcurrent = stdext.Max(cfg.MaxConcurrent, 1)
	cfg.Rate = stdext.Max(cfg.Rate, time.Nanosecond)

	queue := &RateLimitedQueue[T]{
		inChannel:  make(chan T, cfg.BufferSize),
		outChannel: make(chan T, cfg.MaxConcurrent-1),
	}
	go queue.run(cfg.Rate)
	return queue
}

func (r *RateLimitedQueue[T]) run(rate time.Duration) {
	ticker := time.NewTicker(rate)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case t, ok := <-r.inChannel:
			if ok {
				r.outChannel <- t
			} else {
				close(r.outChannel)
				ticker.Stop()
			}
		default:
		}
	}
}

// Push pushes a message onto the queue if the queue has space available within the given timeout.
// Otherwise, it returns a RateLimitedQueueFull error
func (r *RateLimitedQueue[T]) Push(i T, timeout time.Duration) error {
	after := time.After(timeout)
	select {
	case r.inChannel <- i:
		return nil
	case <-after:
		return &RateLimitedQueueFull{}
	}
}

// PushBlocking pushes a message onto the queue, blocking until the queue has space available.
func (r *RateLimitedQueue[T]) PushBlocking(i T) {
	r.inChannel <- i
}

// Pop pops a message off the queue,
// returning the zero value of T and false if no message becomes available (or if the queue is closed) during the timeout.
// The last bool value indicates whether the queue is still open.
func (r *RateLimitedQueue[T]) Pop(timeout time.Duration) (T, bool, bool) {
	after := time.After(timeout)
	select {
	case t, ok := <-r.outChannel:
		return t, ok, ok
	case <-after:
		var nothing T
		return nothing, false, true
	}
}

// PopBlocking pops a message off the queue, but blocks until a message is available or the queue is closed.
// The bool value indicates whether the queue is still open.
func (r *RateLimitedQueue[T]) PopBlocking() (T, bool) {
	v, ok := <-r.outChannel
	return v, ok
}

// PopMultiple waits the duration of the timeout, and returns all the messages that was available in the timeout.
// The bool value indicates whether the queue is still open.
func (r *RateLimitedQueue[T]) PopMultiple(timeout time.Duration) ([]T, bool) {
	after := time.After(timeout)
	available := make([]T, 0)
	open := false
	for {
		select {
		case t, ok := <-r.outChannel:
			open = ok
			if ok {
				available = append(available, t)
			}
		case <-after:
			return available, open
		}
	}
}

// Close starts closing the queue. Pushing any messages onto the queue after calling Close, will panic.
//
// Messages can still be popped until there are none left.
// The queue will only be considered closed once all the messages have been popped.
func (r *RateLimitedQueue[T]) Close() {
	close(r.inChannel)
}

// Consume executes the given function for every message in the queue as it is made available.
// Consume blocks until the queue is closed with Close and all messages have been popped.
func (r *RateLimitedQueue[T]) Consume(f func(T)) {
	for t := range r.outChannel {
		f(t)
	}
}
