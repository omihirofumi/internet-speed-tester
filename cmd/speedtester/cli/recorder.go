package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type recorder struct {
	byteLen int64
	start   time.Time
	lapch   chan Lap
}

func newRecorder(start time.Time, cpun int) *recorder {
	return &recorder{
		start: start,
		lapch: make(chan Lap, cpun),
	}
}

func (r *recorder) download(ctx context.Context, url string, size int) error {
	url = fmt.Sprintf("%s?size=%s", url, size)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	// start recording

	return nil
}
