package main

type Lap struct {
	Bps   float64
	Bytes int64
	delta float64
}

func newLap(byteLen int64, delta float64) Lap {
	var bytes float64
	if delta > 0 {
		bytes = float64(byteLen) / delta
	}
	bps := bytes * 8
	return Lap{
		Bytes: byteLen,
		Bps:   bps,
		delta: delta,
	}
}
