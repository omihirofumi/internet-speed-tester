package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
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

type measureProxy struct {
	io.Reader
	*recorder
}

func (r *recorder) newMeasureProxy(ctx context.Context, reader io.Reader) io.Reader {
	rp := &measureProxy{
		Reader:   reader,
		recorder: r,
	}
	go rp.Watch(ctx, r.lapch)
	return rp
}

func (m *measureProxy) Watch(ctx context.Context, send chan<- Lap) {
	t := time.NewTicker(150 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			byteLen := atomic.LoadInt64(&m.byteLen)
			delta := time.Now().Sub(m.start).Seconds()
			send <- newLap(byteLen, delta)
		case <-ctx.Done():
			return
		}
	}
}
