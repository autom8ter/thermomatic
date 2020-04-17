package client

import (
	"context"
	"log"
	"net"
)

//Cache can persist,fetch, and delete each client's last reading in memory
type Cache interface {
	SetReading(imei uint64, reading *Reading)
	GetReading(imei uint64) (*Reading, bool)
	DeleteReading(imei uint64)
}

//ClientConn represents a single Thermomatic client connection
type ClientConn interface {
	GetConn() net.Conn
	SetIMEI(code uint64)
	GetIMEI() uint64
	GetManager() Manager
	Connect(ctx context.Context)
	Close()
}

type ClientHub interface {
	AddClient(c ClientConn)
	RemoveClient(imei uint64)
}

type Logger interface {
	GetClientLogger() *log.Logger
	GetServerLogger() *log.Logger
}

type Manager interface {
	Logger
	ClientHub
	Cache
}
