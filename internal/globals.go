package internal

import "net"

var (
	IOHandles   = NewSyncMap[int32, int32]()
	Connections = NewSyncMap[int32, net.Conn]()
)
