package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync/atomic"
	"time"
)

var (
	numbersGenerated uint32
	writesFailed uint32
)

func main() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 5; i++ {
		go newClientGeneratingData(i)
	}

	time.Sleep(30 * time.Second)

	avg := numbersGenerated /3
	_, _ = fmt.Fprintf(
		os.Stdout,
		"Total numbers generated: %d\nTotal writes failed: %d\nAverage in 10 sec: %d\n",
		numbersGenerated,
		writesFailed,
		avg,
	)
	if avg < 2_000_000 {
		fmt.Fprintln(os.Stderr, "Requirement (2.000.000 avg per 10 secs) not met")
		os.Exit(1)
	}
}

func newClientGeneratingData(i int) {
	conn, _ := net.Dial("tcp", ":4000")

	for {
		_, err := conn.Write(generateNumber())
		if err != nil {
			atomic.AddUint32(&writesFailed, 1)
		}
	}
}

func generateNumber() []byte {
	number := fmt.Sprintf("%09d\n", rand.Intn(999_999_999))
	atomic.AddUint32(&numbersGenerated, 1)
	return []byte(number)
}
