package gocontainers

import (
	"errors"
	"fmt"
	"sync"
)

type RingBuffer struct {
	buffer []byte
	size   int
	r      int // current read position
	w      int // current write position
	full   bool
	sync.Mutex
}

func NewRingBuffer(size int) *RingBuffer {
	if size <= 0 {
		return nil
	}

	r := &RingBuffer{
		buffer: make([]byte, size),
		size:   size,
		r:      0,
		w:      0,
	}

	return r
}

func (r *RingBuffer) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, errors.New("invalid params")
	}

	if r.WriteableLen() <= 0 {
		return 0, errors.New("buffer is full")
	}

	n := len(p)
	var availLen int
	if r.w >= r.r {
		availLen = (r.size - r.w) + r.r
	} else {
		availLen = r.r - r.w
	}

	if availLen < n {
		p = p[:availLen]
		n = len(p)
	}

	if r.w >= r.r {
		end := r.size - r.w
		if end >= n {
			copy(r.buffer[r.w:], p)
			r.w += n
		} else {
			// p:|--------------|-----------|
			//   |------end-----|--surplus--|
			copy(r.buffer[r.w:], p[:end])
			// write
			surplus := n - end
			copy(r.buffer[0:], p[end:end+surplus])
			r.w = surplus
		}
	} else {
		end := r.r - r.w
		copy(r.buffer[r.w:], p[:end])
		r.w += end
	}

	// ring back
	if r.w == r.size {
		r.w = 0
	}

	if r.w == r.r {
		r.full = true
	}

	return n, nil
}

func (r *RingBuffer) WriteByte(b byte) error {
	if r.WriteableLen() <= 0 {
		return errors.New("buffer is full")
	}

	r.buffer[r.w] = b
	r.w += 1

	if r.size == r.w {
		r.w = 0
	}

	if r.w == r.r {
		r.full = true
	}

	return nil
}

func (r *RingBuffer) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	if r.ReadableLen() <= 0 {
		return 0, errors.New("buffer is empty")
	}

	n := len(p)
	// r is not complete ring
	if r.w > r.r {
		read := r.w - r.r
		if n > read {
			n = read
		}
		copy(p[0:], r.buffer[r.r:r.r+n])
		r.r += n
		return n, nil
	}

	read := (r.size - r.r) + r.w
	if n > read {
		n = read
	}

	if r.r+n <= r.size {
		copy(p[0:], r.buffer[r.r:r.r+n])
	} else {
		read1 := r.size - r.r
		copy(p[0:], r.buffer[r.r:r.size])
		read2 := n - read1
		copy(p[read1:], r.buffer[0:read2])
	}

	r.r = (r.r + n) % r.size
	r.full = false
	return n, nil
}

// just for debug
func (r *RingBuffer) String() string {
	return fmt.Sprintf("Read: %v, Write: %v, Size: %v, Full: %v",
		r.r, r.w, r.size, r.full)
}

func (r *RingBuffer) Reset() {
	r.r = 0
	r.w = 0
	r.full = false
}

func (r *RingBuffer) ReadableLen() int {
	if r.r == r.w && !r.full {
		return 0
	}

	return (r.size - r.r) + r.w
}
func (r *RingBuffer) WriteableLen() int {
	if r.r == r.w && r.full {
		return 0
	}

	return (r.size - r.w) + r.r
}
