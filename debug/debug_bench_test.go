// Copyright (c) 2016 DDN. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
 * Quick benchmark used for figuring out the lowest-overhead method of
 * toggling debug. On a 2016 MBP, the benchmark shows these results:
 *
BenchmarkMutexBool-8 	10000000	       185 ns/op	       0 B/op	       0 allocs/op
BenchmarkAtomicBool-8	2000000000	         1.41 ns/op	       0 B/op	       0 allocs/op
*/

package debug_test

import (
	"sync"
	"sync/atomic"
	"testing"
)

type boolStruct struct {
	sync.Mutex
	enabled bool
}

func (bs *boolStruct) Enabled() bool {
	bs.Lock()
	defer bs.Unlock()
	return bs.enabled
}

func BenchmarkMutexBool(b *testing.B) {
	bs := &boolStruct{}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			bs.Enabled()
		}
	})
}

type atomicStruct struct {
	enabled int32
}

func (as *atomicStruct) Enabled() bool {
	return atomic.LoadInt32(&as.enabled) == 1
}

func BenchmarkAtomicBool(b *testing.B) {
	as := &atomicStruct{}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			as.Enabled()
		}
	})
}
