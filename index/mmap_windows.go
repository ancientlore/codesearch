// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"log"
	"os"
	"syscall"
	"unsafe"
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
	if size == 0 {
		return mmapData{f, nil, 0}
	}
	h, err := syscall.CreateFileMapping(syscall.Handle(f.Fd()), nil, syscall.PAGE_READONLY, uint32(size>>32), uint32(size), nil)
	if err != nil {
		f.Close()
		log.Fatalf("CreateFileMapping %s: %v", n, err)
	}

	addr, err := syscall.MapViewOfFile(h, syscall.FILE_MAP_READ, 0, 0, 0)
	if err != nil {
		f.Close()
		log.Fatalf("MapViewOfFile %s: %v", n, err)
	}
	data := (*[1 << 30]byte)(unsafe.Pointer(addr))
	return mmapData{f, data[:size], uintptr(h)}
}

func (data mmapData) close() {
	var err1, err2 error
	n := data.f.Name()
	if data.h != 0 {
		err1 = syscall.Close(syscall.Handle(data.h))
	}
	if data.f != nil {
		err2 = data.f.Close()
	}
	if err1 != nil {
		log.Fatalf("close handle %s: %v", n, err1)
	}
	if err2 != nil {
		log.Fatalf("close %s: %v", n, err2)
	}
}
