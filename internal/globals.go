package internal

import "net"

var (
	Connections = NewSyncMap[int32, net.Conn]()
	IOHandles   = NewSyncMap[int32, int32]()
)
