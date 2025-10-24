// SPDX-License-Identifier: BSD-2-Clause
//
// Copyright (c) Lewis Cook <lcook@FreeBSD.org>
package cache

import (
	"sync/atomic"
	"unsafe"
)

type RingBuffer[T any] struct {
	buffer []unsafe.Pointer
	size   uint64
	write  atomic.Uint64
	read   atomic.Uint64
}

func NewRingBuffer[T any](size uint64) *RingBuffer[T] {
	return &RingBuffer[T]{
		buffer: make([]unsafe.Pointer, size),
		size:   size,
	}
}

func (r *RingBuffer[T]) Size() uint64 {
	return r.write.Load() - r.read.Load()
}

func (r *RingBuffer[T]) Add(value T) {
	if r.Size() >= r.size {
		r.read.Add(1)
	}

	ptr := unsafe.Pointer(&value)
	idx := r.write.Load() % r.size

	atomic.StorePointer(&r.buffer[idx], ptr)

	r.write.Add(1)
}

func (r *RingBuffer[T]) Slice() []T {
	var result []T

	for i := range r.Size() {
		idx := (r.read.Load() + i) % r.size
		ptr := atomic.LoadPointer(&r.buffer[idx])

		if ptr != nil {
			result = append(result, *(*T)(ptr))
		}
	}

	return result
}

func (r *RingBuffer[T]) ForEach(fn func(*T)) {
	for i := range r.Size() {
		idx := (r.read.Load() + i) % r.size
		ptr := atomic.LoadPointer(&r.buffer[idx])

		if ptr != nil {
			fn((*T)(ptr))
		}
	}
}
