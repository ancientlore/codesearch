// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"log"
	"os"
	"syscall"
)

func mmapFile(f *os.File) mmapData {
	st, err := f.Stat()
	n := f.Name()
	if err != nil {
		f.Close()
		log.Fatal(err)
	}
	size := st.Size()
	if int64(int(size+4095)) != size+4095 {
		f.Close()
		log.Fatalf("%s: too large for mmap", n)
	}
	n := int(size)
	if n == 0 {
		return mmapData{f, nil, 0}
	}
	data, err := syscall.Mmap(int(f.Fd()), 0, (n+4095)&^4095, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		f.Close()
		log.Fatalf("mmap %s: %v", n, err)
	}
	return mmapData{f, data[:n], 0}
}

func (data mmapData) close() {
	if data.f != nil {
		n := data.f.Name()
		err := data.f.Close()
		if err != nil {
			log.Fatalf("close %s: %v", n, err)
		}
	}
}
