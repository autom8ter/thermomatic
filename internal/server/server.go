package server

import (
	"context"
	"fmt"
	"github.com/autom8ter/thermomatic/internal/client"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

//Config holds the configuration requirements to start the server
type Config struct {
	TcpPort         int
	HttpPort        int
	ClientLogPrefix string
	ServerLogPrefix string
}

//Server serves tcp connections for logging iot device readings and serves http endpoints for iot reading statistics/analysis
type Server struct {
	tcpLis    *net.TCPListener
	httpPort  string
	mux       *http.ServeMux
	serverLog *log.Logger
	clientLog *log.Logger
	wg        *sync.WaitGroup
	clientMu  *sync.Mutex
	clients   map[uint64]client.ClientConn
	readings  map[uint64]*client.Reading
	readingMu *sync.Mutex
}

//NewServer creates a new server instance from the given config
func NewServer(config *Config) (*Server, error) {
	serverLog := log.New(os.Stdin, config.ServerLogPrefix, log.LstdFlags)
	clientLog := log.New(os.Stderr, config.ClientLogPrefix, log.LstdFlags)
	tcpLis, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: config.TcpPort,
	})
	if err != nil {
		return nil, err
	}
	return &Server{
		tcpLis:    tcpLis,
		httpPort:  fmt.Sprintf(":%v", config.HttpPort),
		mux:       http.NewServeMux(),
		serverLog: serverLog,
		clientLog: clientLog,
		clientMu:  &sync.Mutex{},
		wg:        &sync.WaitGroup{},
		clients:   map[uint64]client.ClientConn{},
		readingMu: &sync.Mutex{},
		readings:  map[uint64]*client.Reading{},
	}, nil
}

// Listen starts the tcp and http server
func (s *Server) Listen(ctx context.Context) {
	s.setupRoutes()
	defer s.tcpLis.Close()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.serverLog.Println("starting tcp server!")
		for {
			select {
			case <-ctx.Done():
				break
			default:
				if err := s.tcpLis.SetDeadline(time.Now().Add(1 * time.Minute)); err != nil {
					s.serverLog.Printf("[ERROR] failed to accept tcp connection: %s", err.Error())
					continue
				}
				conn, err := s.tcpLis.Accept()
				if err != nil {
					s.serverLog.Printf("[ERROR] failed to accept tcp connection: %s", err.Error())
					continue
				}
				clientConn, err := client.NewClient(conn, s, s, s.clientLog, s.serverLog)
				if err != nil {
					s.serverLog.Printf("[ERROR] failed to create client: %s", err.Error())
					continue
				}
				s.wg.Add(1)
				go func(conn client.ClientConn) {
					defer s.wg.Done()
					conn.Connect(ctx)
				}(clientConn)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.serverLog.Println("starting http server!")
		if err := http.ListenAndServe(s.httpPort, s.mux); err != nil {
			log.Fatalf("[FATAL] %s", err.Error())
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.wg.Wait()
	}()
	wg.Wait()
}

//AddClient adds a client connection to manage
func (s *Server) AddClient(client client.ClientConn) {
	s.clientMu.Lock()
	defer s.clientMu.Unlock()
	s.clients[client.GetIMEI()] = client
}

//RemoveClient removes the client connection
func (s *Server) RemoveClient(imei uint64) {
	s.clientMu.Lock()
	defer s.clientMu.Unlock()
	if _, ok := s.clients[imei]; ok {
		delete(s.clients, imei)
	}
}

func (c *Server) SetReading(imei uint64, reading *client.Reading) {
	c.readingMu.Lock()
	defer c.readingMu.Unlock()
	c.readings[imei] = reading
}

func (c *Server) GetReading(imei uint64) (*client.Reading, bool) {
	c.readingMu.Lock()
	defer c.readingMu.Unlock()
	if reading, ok := c.readings[imei]; ok {
		return reading, true
	}
	return nil, false
}

func (c *Server) DeleteReading(imei uint64) {
	c.readingMu.Lock()
	defer c.readingMu.Unlock()
	if _, ok := c.readings[imei]; ok {
		delete(c.readings, imei)
	}
}
