package io

import "time"

//go:wasmimport extism:host/user io.ready
func _ready(ioHandle int32) int32

func Ready(ioHandle int32) int32 {
	var result int32
	for {
		result = _ready(ioHandle)
		if result != 0 {
			return result
		}
		time.Sleep(1 * time.Millisecond)
	}
}
