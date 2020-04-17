package client

import (
	"context"
	"fmt"
	"github.com/autom8ter/thermomatic/internal/common"
	"github.com/autom8ter/thermomatic/internal/imei"
	"io"
	"log"
	"net"
	"time"
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
	GetHub() ClientHub
	GeCache() Cache
	Connect(ctx context.Context)
	Close()
}

type ClientHub interface {
	AddClient(c ClientConn)
	RemoveClient(imei uint64)
}

//client implements ClientConn
type client struct {
	conn net.Conn
	//imei is the clients imei(unique identifier)
	imei uint64
	//clientLog logs readings
	clientLog *log.Logger
	//serverLog logs errors and server-related events
	serverLog *log.Logger
	//handleErr handles all errors during the lifecycle of the connection
	handleErr func(c ClientConn, err error)
	//handleReading handles all client readings during the lifecycle of the connection
	handleReading func(c ClientConn, reading *Reading) error
	//handleLogin handles the clients first message to log them in. The connection will be closed if an error is returned
	handleLogin func(c ClientConn) error
	//handleDone is executed when the client connection is closing
	handleDone func(c ClientConn)
	close      chan struct{}
	hub        ClientHub
	//cache is used to persist the clients readings
	cache Cache
}

//NewClient creates a new ClientConn with default event handlers. clientLog will be used to log readings
func NewClient(conn net.Conn, hub ClientHub, cache Cache, clientLog, serverLog *log.Logger) (ClientConn, error) {
	client := &client{
		conn:      conn,
		clientLog: clientLog,
		serverLog: serverLog,
		handleErr: func(c ClientConn, err error) {
			serverLog.Printf("[ERROR] %v error: %s", c.GetIMEI(), err)
		},
		hub:   hub,
		close: make(chan struct{}, 1),
		cache: cache,
	}
	client.handleLogin = func(c ClientConn) error {
		if err := conn.SetReadDeadline(time.Now().Add(1 * time.Second)); err != nil {
			return err
		}
		b := make([]byte, 15) //read imei from connection
		if _, err := conn.Read(b); err != nil {
			return err
		}
		//if _, err := io.ReadFull(conn, b); err != nil {
		//	return err
		//}
		code, err := imei.Decode(b)
		if err != nil {
			return err
		}
		c.SetIMEI(code)
		hub.AddClient(c)
		return nil
	}
	client.handleReading = func(c ClientConn, message *Reading) error {
		if c.GetIMEI() == 0 {
			return fmt.Errorf("failed handle reading: empty imei code")
		}
		message.Log(c.GetIMEI(), clientLog)
		client.cache.SetReading(c.GetIMEI(), message)
		return nil
	}
	client.handleDone = func(c ClientConn) {
		client.cache.DeleteReading(c.GetIMEI())
		client.hub.RemoveClient(c.GetIMEI())
	}
	return client, nil
}

//Connect handles the lifecycle of the client connection using the clients event handlers(see other methods to override)
func (c *client) Connect(ctx context.Context) {
	defer c.conn.Close()
	for {
		select {
		default:
			if c.GetIMEI() == 0 {
				//handleLogin when the client first establishes a conectionn
				if err := c.handleLogin(c); err != nil {
					c.handleErr(c, fmt.Errorf("client login: %s", err))
					c.Close()
					return
				}

			}
			b := make([]byte, 40) //read imei from connection
			if _, err := c.GetConn().Read(b); err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					c.handleErr(c, fmt.Errorf("client timeout: %s", err))
					c.Close()
					return
				}
				if err == io.EOF {
					return
				}
				c.handleErr(c, fmt.Errorf("failed to read message: %s", err))
				return
			}
			var reading = new(Reading)
			if len(b) >= common.MinReadingLength {
				ok, err := reading.Decode(b)
				if err != nil {
					c.handleErr(c, fmt.Errorf("decode reading: %s", err))
					continue
				}
				if ok {
					if err := c.handleReading(c, reading); err != nil {
						c.handleErr(c, fmt.Errorf("handle reading: %s", err))
					}
				}
			}
		case <-ctx.Done():
			c.handleDone(c)
			break
		case <-c.close:
			c.handleDone(c)
			break
		}
	}
}

//GetConn gets the clients connection
func (c *client) GetConn() net.Conn {
	return c.conn
}

//SetIMEI sets the clients imei code
func (c *client) SetIMEI(code uint64) {
	c.imei = code
}

//GetIMEI retrieves the clients imei code
func (c *client) GetIMEI() uint64 {
	return c.imei
}

func (c *client) GetHub() ClientHub {
	return c.hub
}

func (c *client) GeCache() Cache {
	return c.cache
}

//Close is used to close a client connection
func (c *client) Close() {
	c.close <- struct{}{}
}

//OnLogin is used to change the default behavior of a client on login
func (c *client) OnLogin(handler func(c ClientConn, err error)) {
	c.handleErr = handler
}

//OnErr is used to change the default behavior of a client whenever an error occurs during the lifecycle of its connection
func (c *client) OnErr(handler func(c ClientConn, err error)) {
	c.handleErr = handler
}

//OnRead is used to change the default behavior when a client sends a reading
func (s *client) OnRead(handler func(c ClientConn, reading *Reading) error) {
	s.handleReading = handler
}

//OnDone is used to change the default behavior when a client is closing
func (c *client) OnDone(handler func(c ClientConn)) {
	c.handleDone = handler
}
