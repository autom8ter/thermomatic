package client

import (
	"context"
	"fmt"
	"github.com/autom8ter/thermomatic/internal/common"
	"github.com/autom8ter/thermomatic/internal/imei"
	"io"
	"net"
	"time"
)

//client implements ClientConn
type client struct {
	conn net.Conn
	//imei is the clients imei(unique identifier)
	imei    uint64
	manager Manager
	//handleErr handles all errors during the lifecycle of the connection
	handleErr func(c ClientConn, err error)
	//handleReading handles all client readings during the lifecycle of the connection
	handleReading func(c ClientConn, reading *Reading) error
	//handleLogin handles the clients first message to log them in. The connection will be closed if an error is returned
	handleLogin func(c ClientConn) error
	//handleDone is executed when the client connection is closing
	handleDone func(c ClientConn)
	close      chan struct{}
}

//NewClient creates a new ClientConn with default event handlers. clientLog will be used to log readings
func NewClient(conn net.Conn, manager Manager) (ClientConn, error) {
	client := &client{
		conn:    conn,
		manager: manager,
		handleErr: func(c ClientConn, err error) {
			manager.GetServerLogger().Printf("[ERROR] %v error: %s", c.GetIMEI(), err)
		},
		close: make(chan struct{}, 1),
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
		c.GetManager().AddClient(c)
		return nil
	}
	client.handleReading = func(c ClientConn, message *Reading) error {
		if c.GetIMEI() == 0 {
			return fmt.Errorf("failed handle reading: empty imei code")
		}
		message.Log(c.GetIMEI(), c.GetManager().GetClientLogger())
		c.GetManager().SetReading(c.GetIMEI(), message)
		return nil
	}
	client.handleDone = func(c ClientConn) {
		c.GetManager().DeleteReading(c.GetIMEI())
		c.GetManager().RemoveClient(c.GetIMEI())
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
			if err := c.GetConn().SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
				c.handleErr(c, fmt.Errorf("client read timeout: %s", err))
				c.Close()
				return
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

func (c *client) GetManager() Manager {
	return c.manager
}

//Close is used to close a client connection
func (c *client) Close() {
	c.close <- struct{}{}
}
