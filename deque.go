package schemagen

import (
	"errors"
)

type Deque struct {
	buf  []rune
	left int
	next int
	cnt  int
}

func NewDeque() *Deque {
	it := new(Deque)
	initSz := 1 << 10
	it.buf = make([]rune, initSz, initSz)
	return it
}

func (this *Deque) expand() {
	newSize := 2 * cap(this.buf)
	newBuf := make([]rune, newSize, newSize)
	left := this.left
	copy(newBuf, this.buf[left:])
	if this.left > 0 {
		copy(newBuf[cap(this.buf)-left:], this.buf[:left])
	}
	this.left = 0
	this.next = cap(this.buf)
	this.buf = newBuf
}

func (this *Deque) isEmpty() bool {
	return this.cnt == 0
}

func (this *Deque) inc(n int) int {
	return (n + 1) % cap(this.buf)
}

func (this *Deque) dec(n int) int {
	return (n - 1) % cap(this.buf)
}

func (this *Deque) pushR(r rune) {
	if this.cnt == cap(this.buf) {
		this.expand()
	}
	if this.cnt == 0 {
		this.buf[0] = r
		this.left = 0
		this.next = 1
	} else {
		this.buf[this.next] = r
		this.next = this.inc(this.next)
	}
	this.cnt++
}

func (this *Deque) popR() (rune, error) {
	if this.cnt == 0 {
		return 0, errors.New("pop from empty deque")
	}
	r := this.buf[this.dec(this.next)]
	this.left = this.dec(this.next)
	this.cnt--
	return r, nil
}

func (this *Deque) pushL(r rune) {
	if this.cnt == cap(this.buf) {
		this.expand()
	}
	if this.cnt == 0 {
		this.buf[0] = r
		this.left = 0
		this.next = 1
	} else {
		this.buf[this.dec(this.left)] = r
		this.left = this.dec(this.left)
	}
	this.cnt++
}

func (this *Deque) popL() (rune, error) {
	if this.cnt == 0 {
		return 0, errors.New("pop from empty deque")
	}
	r := this.buf[this.left]
	this.left = this.inc(this.left)
	this.cnt--
	return r, nil
}
