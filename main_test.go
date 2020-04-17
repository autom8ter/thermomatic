package main_test

import (
	"encoding/json"
	"github.com/autom8ter/thermomatic/internal/client"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	url2 "net/url"
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
		//go get the latest value from endpoint
		target, err := url2.Parse("http://localhost:1338/readings")
		if err != nil {
			t.Fatal(err.Error())
		}
		q := target.Query()
		q.Set("imei", imei)
		target.RawQuery = q.Encode()
		resp, err := http.DefaultClient.Get(target.String())
		if err != nil {
			t.Fatal(err.Error())
		}
		defer resp.Body.Close()
		reading := &client.Reading{}
		bits, _ = ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(bits, reading); err != nil {
			t.Fatal(err.Error())
		}
	}
}

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
	t.Logf("TestE2ELoad load: %v time: %vms", load, time.Since(now).Nanoseconds()/1000000)
}
