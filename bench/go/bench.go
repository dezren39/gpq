package main

import (
	"log"
	"math/rand"
	_ "net/http/pprof"
	"time"

	"github.com/JustinTimperio/gpq"
)

var (
	total      int  = 10000000
	syncToDisk bool = false
	print      bool = false
	maxBuckets int  = 10

	sent         uint64
	received     uint64
	missed       int64
	hits         int64
	lastPriority int64
)

func main() {

	queue, err := gpq.NewGPQ[int](maxBuckets, syncToDisk, "/tmp/gpq/")
	if err != nil {
		log.Fatalln(err)
	}

	timer := time.Now()
	for i := 0; i < total; i++ {
		p := rand.Intn(maxBuckets)
		timer := time.Now()
		err := queue.EnQueue(
			i,
			int64(p),
			false,
			time.Minute,
			false,
			10*time.Minute,
		)
		if err != nil {
			log.Fatalln(err)
		}
		if print {
			log.Println("EnQueue", p, time.Since(timer))
		}
		sent++
	}

	for total > int(received) {
		timer := time.Now()
		priority, item, err := queue.DeQueue()
		if err != nil {
			log.Println(err)
			lastPriority = 0
			continue
		}
		received++
		if print {
			log.Println("DeQueue", priority, received, item, time.Since(timer))
		}

		if lastPriority > priority {
			missed++
		} else {
			hits++
		}
		lastPriority = priority
	}

	log.Println("Sent", sent, "Received", received, "Finished in", time.Since(timer), "Missed", missed, "Hits", hits)

}
