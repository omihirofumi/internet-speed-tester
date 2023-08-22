package main

import (
	"context"
	"log"
	"runtime"
	"time"
)

const (
	endpoint = "http://localhost"
	size     = 1000
)

func main() {
	r := newRecorder(time.Now(), runtime.NumCPU())
	if err := r.download(context.Background(), endpoint, size); err != nil {
		log.Fatal(err)
	}
}
