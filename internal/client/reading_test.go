package client_test

import (
	"github.com/autom8ter/thermomatic/internal/client"
	"log"
	"testing"
)

func init() {
	r := &client.Reading{
		Temperature:  102.45,
		Altitude:     0,
		Latitude:     40.936976099999995,
		Longitude:    -165.00857800000001,
		BatteryLevel: .12,
	}
	bits, err := r.Encode()
	if err != nil {
		log.Fatal(err.Error())
	}
	singleEncodedReading = bits
}

var singleEncodedReading []byte

//TestDecode decodes and validates each reading
func TestDecode(t *testing.T) {
	tests := []struct {
		Name    string
		Reading *client.Reading
		Pass    bool
	}{
		{
			Name: "accurate reading (1)",
			Reading: &client.Reading{
				Temperature:  50.77,
				Altitude:     5280,
				Latitude:     39.936976099999995,
				Longitude:    -105.00857800000001,
				BatteryLevel: 95,
			},
			Pass: true,
		},
		{
			Name: "accurate reading (2)",
			Reading: &client.Reading{
				Temperature:  102.45,
				Altitude:     0,
				Latitude:     40.936976099999995,
				Longitude:    -165.00857800000001,
				BatteryLevel: .12,
			},
			Pass: true,
		},
		{
			Name: "inacurate reading (1) - temperature",
			Reading: &client.Reading{
				Temperature:  5000, //too hot
				Altitude:     5280,
				Latitude:     39.936976099999995,
				Longitude:    -105.00857800000001,
				BatteryLevel: 95,
			},
			Pass: false,
		},
		{
			Name: "inacurate reading (2) - altitude",
			Reading: &client.Reading{
				Temperature:  32,
				Altitude:     500000, //too high
				Latitude:     39.936976099999995,
				Longitude:    -105.00857800000001,
				BatteryLevel: 95,
			},
			Pass: false,
		},
		{
			Name: "inacurate reading (3) - battery",
			Reading: &client.Reading{
				Temperature:  50,
				Altitude:     5280,
				Latitude:     39.936976099999995,
				Longitude:    -105.00857800000001,
				BatteryLevel: 195, //too high
			},
			Pass: false,
		},
		{
			Name: "inacurate reading (3) - latitude",
			Reading: &client.Reading{
				Temperature:  50,
				Altitude:     5280,
				Latitude:     99.936976099999995, //should be under 90
				Longitude:    -105.00857800000001,
				BatteryLevel: 55,
			},
			Pass: false,
		},
		{
			Name: "inacurate reading (4) - longitude",
			Reading: &client.Reading{
				Temperature:  50,
				Altitude:     5280,
				Latitude:     39.936976099999995,
				Longitude:    -199.00857800000001, //should be under 180
				BatteryLevel: 55,
			},
			Pass: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			b, err := test.Reading.Encode()
			if err != nil {
				t.Errorf("error : %s", err)
			}
			if len(b) != 40 {
				t.Errorf("expected 40 lenght byteslice actual = %v", len(b))
			}
			reading := &client.Reading{}
			ok, err := reading.Decode(b)
			if err != nil && test.Pass {
				t.Errorf("error : %s reading : %s", err.Error(), test.Reading.String(0))
			}
			if !ok && test.Pass {
				t.Errorf("failed to decode reading: %s", test.Reading.String(0))
			}
		})
	}
}

//go test -v -bench=.
//BenchmarkDecode-12     20000000                75.4 ns/op             0 B/op          0 allocs/op
func BenchmarkDecode(b *testing.B) {
	b.ReportAllocs()
	r := new(client.Reading)
	for i := 0; i < b.N; i++ {
		r.Decode(singleEncodedReading)
	}
}
