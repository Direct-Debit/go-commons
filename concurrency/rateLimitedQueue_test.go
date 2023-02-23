package concurrency

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
	"testing"
	"time"
)

type queueMeasurement struct {
	id      string
	created time.Time
	done    time.Time
	success bool
}

func TestRateLimitedQueue(t *testing.T) {
	queue := NewRateLimitedQueue[queueMeasurement](RateLimitedQueueConfig{
		BufferSize:    16,
		MaxConcurrent: 1,
		Rate:          time.Millisecond * 50,
	})

	results := make([]queueMeasurement, 0)

	var pops sync.WaitGroup
	pops.Add(1)
	go func() {
		queue.Consume(func(m queueMeasurement) {
			m.done = time.Now()
			m.success = true
			log.Infof("Pop: queue item %s popped at %v", m.id, m.done)
			log.Infof("Queue State: %+v", queue.State())
			results = append(results, m)
		})
		pops.Done()
	}()

	var pushes sync.WaitGroup
	pushes.Add(1)
	go func() {
		defer pushes.Done()
		for i := 0; i <= 17; i++ {
			m, ok := pushQueueItem(queue, fmt.Sprintf("init.%d", i))
			if !ok {
				results = append(results, m)
			}
		}
	}()

	for i := 0; i <= 128; i++ {
		time.Sleep(time.Millisecond * 150)
		pushes.Add(1)
		go func(ii int) {
			defer pushes.Done()

			pushCount := 1
			if (ii+1)%64 == 0 {
				pushCount = 64
			}
			for j := 0; j <= pushCount; j++ {
				m, ok := pushQueueItem(queue, fmt.Sprintf("%d.%d", ii, j))
				if !ok {
					results = append(results, m)
				}
			}
		}(i)
	}

	log.Infof("Waiting for pushes to finish...")
	pushes.Wait()
	queue.Close()
	log.Infof("Waiting for pops to finish...")
	pops.Wait()

	for _, r := range results {
		fmt.Printf("%s,%v,%s,%s\n", r.id, r.success, r.created.Format(time.RFC3339Nano), r.done.Format(time.RFC3339Nano))
	}
}

func pushQueueItem(q *RateLimitedQueue[queueMeasurement], id string) (queueMeasurement, bool) {
	item := queueMeasurement{
		id:      id,
		created: time.Now(),
	}
	err := q.Push(item, time.Millisecond*150)
	if err != nil {
		item.done = time.Now()
		item.success = false
		log.Infof("Failed: queue item %s timed out at at %v", item.id, item.done)
	} else {
		log.Infof("Push: queue item %s pushed at %v", item.id, item.created)
	}
	log.Infof("Queue State: %+v", q.State())
	return item, err == nil
}
