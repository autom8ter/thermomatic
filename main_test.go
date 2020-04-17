package main_test

import (
	"github.com/autom8ter/thermomatic/internal/client"
	"log"
	"net"
	"testing"
)

var (
	imei = "450154603277518"
)

func TestE2E(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:1337")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer conn.Close()
	//login
	if _, err := conn.Write([]byte(imei)); err != nil {
		log.Fatal(err.Error())
	}
	var tests = []struct {
		Reading *client.Reading
	}{
		{
			Reading: &client.Reading{
				Temperature:  150,
				Altitude:     5280,
				Latitude:     39.936976099999995,
				Longitude:    -105.00857800000001,
				BatteryLevel: 50,
			},
		},
		{
			Reading: &client.Reading{
				Temperature:  150,
				Altitude:     5280,
				Latitude:     39.936976099999995,
				Longitude:    -105.00857800000001,
				BatteryLevel: 50,
			},
		},
		{
			Reading: &client.Reading{
				Temperature:  150,
				Altitude:     5280,
				Latitude:     39.936976099999995,
				Longitude:    -105.00857800000001,
				BatteryLevel: 50,
			},
		},
		{
			Reading: &client.Reading{
				Temperature:  150,
				Altitude:     5280,
				Latitude:     39.936976099999995,
				Longitude:    -105.00857800000001,
				BatteryLevel: 50,
			},
		},
	}
	for _, test := range tests {
		bits, err := test.Reading.Encode()
		if err != nil {
			t.Fatal(err.Error())
		}
		if _, err := conn.Write(bits); err != nil {
			t.Fatal(err.Error())
		}
	}
}
