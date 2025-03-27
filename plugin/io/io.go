package io

import "time"

//go:wasmimport wape:host/env io.ready
func _ready(handle int32) int32

func Ready(handle int32) int32 {
	var result int32
	for {
		result = _ready(handle)
		if result != 0 {
			return result
		}
		time.Sleep(1 * time.Millisecond)
	}
}
