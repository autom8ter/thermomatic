// Package common implements utilities & functionality commonly consumed by the
// rest of the packages.
package common

import (
	"fmt"
)

type ErrType string

const (
	// ErrNotImplemented is raised throughout the codebase of the challenge to
	// denote implementations to be done by the candidate.
	ErrNotImplemented ErrType = "not implemented"
	ErrInvalidImei    ErrType = "imei: imei: invalid"
	ErrChecksum       ErrType = "imei: invalid checksum"
	ErrReadingBytes   ErrType = "reading isn't at least 40 bytes long."
	ErrReadingTemp    ErrType = "the temperature reading of the device is invalid. Celcius. Min/Max: [-300, 300]"
	ErrReadingAlt     ErrType = "the altitude reading of the device is invalid. Meters. Min/Max: [-20000, 20000]"
	ErrReadingLat     ErrType = "the latitude reading of the device is invalid. Degrees. Min/Max: [-90, 90]"
	ErrReadingLon     ErrType = "the longitude reading of the device is invalid. Degrees. Min/Max: [-180, 180]"
	ErrReadingBattery ErrType = "the battery level of the device is invalid. Percentage. Min/Max: (0, 100]"
)

const (
	ASCIIZero        = 48
	MinImeiLength    = 15
	MinReadingLength = 40
)

func Wrap(typ ErrType, details string) error {
	return fmt.Errorf("%s  - %s", typ, details)
}

type Stats struct {
	GoRoutines        int    `json:"goroutines"`
	ClientConnections int    `json:"clientConnections"`
	CPUs              int    `json:"cpus"`
	Version           string `json:"version"`
}
