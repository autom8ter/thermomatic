package client

import (
	"encoding/binary"
	"fmt"
	"github.com/autom8ter/thermomatic/internal/common"
	"math"
	"time"
)

// Reading is the set of device readings.
type Reading struct {
	// Temperature denotes the temperature reading of the message.
	Temperature float64 `json:"temperature"`

	// Altitude denotes the altitude reading of the message.
	Altitude float64 `json:"altitude"`

	// Latitude denotes the latitude reading of the message.
	Latitude float64 `json:"latitude"`

	// Longitude denotes the longitude reading of the message.
	Longitude float64 `json:"longitude"`

	// BatteryLevel denotes the battery level reading of the message.
	BatteryLevel float64 `json:"batteryLevel"`

	Timestamp time.Time `json:"timestamp"`
}

//String returns a human readable string
func (r Reading) String(imei uint64) string {
	return fmt.Sprintf(`%v,%v,%v,%v,%v,%v,%v\n`, time.Now().Unix(), imei, r.Temperature, r.Altitude, r.Latitude, r.Longitude, r.BatteryLevel)

}

//validate returns true with no error if the reading is valid. it also returns an error message if the reading is invalid
func (r Reading) validate() (bool, error) {
	if r.Temperature < -300 || r.Temperature > 300 {
		return false, common.Wrap(common.ErrReadingTemp, fmt.Sprintf("value: %v", r.Temperature))
	}
	if r.BatteryLevel < 0 || r.BatteryLevel > 100 {
		return false, common.Wrap(common.ErrReadingBattery, fmt.Sprintf("value: %v", r.BatteryLevel))
	}
	if r.Altitude < -20000 || r.Altitude > 20000 {
		return false, common.Wrap(common.ErrReadingAlt, fmt.Sprintf("value: %v", r.Altitude))
	}
	if r.Latitude < -90 || r.Latitude > 90 {
		return false, common.Wrap(common.ErrReadingLat, fmt.Sprintf("value: %v", r.Latitude))
	}
	if r.Longitude < -180 || r.Longitude > 180 {
		return false, common.Wrap(common.ErrReadingLon, fmt.Sprintf("value: %v", r.Longitude))
	}
	return true, nil
}

// Decode decodes the reading message payload in the given b into r.
//
// If any of the fields are outside their valid min/max ranges ok will be unset.
//
// Decode does NOT allocate under any condition. Additionally, it panics if b
// isn't at least 40 bytes long.
func (r *Reading) Decode(b []byte) (bool, error) {
	if len(b) < common.MinReadingLength {
		panic(common.Wrap(common.ErrInvalidImei, fmt.Sprintf("invalid imei: %s", string(b))))
	}
	r.Temperature = math.Float64frombits(binary.BigEndian.Uint64(b[0:8]))
	r.Altitude = math.Float64frombits(binary.BigEndian.Uint64(b[8:16]))
	r.Latitude = math.Float64frombits(binary.BigEndian.Uint64(b[16:24]))
	r.Longitude = math.Float64frombits(binary.BigEndian.Uint64(b[24:32]))
	r.BatteryLevel = math.Float64frombits(binary.BigEndian.Uint64(b[32:40]))
	r.Timestamp = time.Now()
	return r.validate()
}

//Log uses the provided logger to log the reading as a human readable string
func (r Reading) Log(code uint64, logger Printer) {
	logger.Printf("record = %s", r.String(code))
}

//Encode encodes a reading to a byteslice
func (r Reading) Encode() ([]byte, error) {
	var (
		body  = make([]byte, 0, 40)
		field = make([]byte, 8)
	)
	binary.BigEndian.PutUint64(field, math.Float64bits(r.Temperature))
	body = append(body, field...)
	binary.BigEndian.PutUint64(field, math.Float64bits(r.Altitude))
	body = append(body, field...)
	binary.BigEndian.PutUint64(field, math.Float64bits(r.Latitude))
	body = append(body, field...)
	binary.BigEndian.PutUint64(field, math.Float64bits(r.Longitude))
	body = append(body, field...)
	binary.BigEndian.PutUint64(field, math.Float64bits(r.BatteryLevel))
	body = append(body, field...)
	return body, nil
}