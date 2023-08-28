package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
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
	proxy := r.newMeasureProxy(ctx, resp.Body)

	if _, err := io.Copy(io.Discard, proxy); err != nil {
		return err
	}

	return nil
}

func (r *recorder) upload(ctx context.Context, url string, size int) error {
	// start measure
	proxy := r.newMeasureProxy(ctx, rand.New(rand.NewSource(0)))
	req, err := http.NewRequest("POST", url, proxy)
	if err != nil {
		return err
	}
	req.ContentLength = int64(size)
	req.Header.Set("Content-Type", "application/octet-stream")
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
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
