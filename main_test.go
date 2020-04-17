package main_test

import (
	"github.com/autom8ter/thermomatic/internal/client"
	"log"
	"net"
	"testing"
	"time"
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
		//resp, err := http.Get(fmt.Sprintf("http://localhost:1337/readings?=%s",imei))
		//if err != nil {
		//	t.Fatal(err.Error())
		//}
		//defer resp.Body.Close()
		//reading := &client.Reading{}
		//if err := json.NewDecoder(resp.Body).Decode(reading); err != nil {
		//	t.Fatal(err.Error())
		//}
	}
	for i := 0; i < 10000; i++ {

	}
}


//go test -v -bench=.
//pkg: github.com/autom8ter/thermomatic/internal/client
//BenchmarkDecode1-12     20000000                75.4 ns/op             0 B/op          0 allocs/op
func TestE2ELoad(t *testing.T) {
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
	}
	const load = 10000
	now := time.Now()
	for i := 0; i < load; i++ {
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
	//57ms on average for 10,000 messages(on averagee)
	t.Logf("TestE2ELoad load: %v time: %vms", load, time.Since(now).Nanoseconds()/ 1000000)
}
