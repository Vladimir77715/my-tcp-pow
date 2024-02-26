package reader

import (
	"errors"
	"io"
	"time"
)

var (
	ErrBuffSizeGraterThanSlice = errors.New("buffer slice grater than expected slice")
	ErrTimeout                 = errors.New("timeout")
)

var DefTimerDuration = 10 * time.Second

type chanStruct struct {
	n int
	e error
}

type WithTimeout struct {
	r io.Reader
	t time.Duration
	s int
}

func New(r io.Reader, t time.Duration, s int) *WithTimeout {
	return &WithTimeout{r: r, t: t, s: s}
}

func (r *WithTimeout) Read(b []byte) (int, error) {
	bChan := make(chan chanStruct)

	go func(b []byte) {
		buff := make([]byte, r.s)
		n, e := r.r.Read(buff)
		if len(buff) > len(b) {
			bChan <- chanStruct{n: 0, e: ErrBuffSizeGraterThanSlice}
			return
		}
		for i := 0; i < len(buff); i++ {
			b[i] = buff[i]
		}

		bChan <- chanStruct{n: n, e: e}
	}(b)

	select {
	case val := <-bChan:
		return val.n, val.e
	case <-time.NewTimer(r.t).C:
		return 0, ErrTimeout
	}
}
