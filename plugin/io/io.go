package io

import "time"

// DefaultDelay is the default delay for [Ready]
var DefaultDelay = time.Millisecond

//go:wasmimport wape:host/env io.ready
func _ready(handle int32) int32

// Ready returns the number of bytes ready to be read/written or negative number in case of error.
// Waits for the [DefaultDelay] between each check of the handle.
func Ready(handle int32) int32 {
	return ReadyWithDelay(handle, DefaultDelay)
}

// ReadyWithDelay returns the number of bytes ready to be read/written or negative number in case of error.
// Waits for the specified delay between each check of the handle.
func ReadyWithDelay(handle int32, delay time.Duration) int32 {
	var result int32
	for {
		result = _ready(handle)
		if result != 0 {
			return result
		}
		time.Sleep(delay)
	}
}
