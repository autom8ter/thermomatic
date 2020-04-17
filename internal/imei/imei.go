// Package imei implements an IMEI decoder.
package imei

// NOTE: for more information about IMEI codes and their structure you may
// consult with:
//
// https://en.wikipedia.org/wiki/International_Mobile_Equipment_Identity.

import (
	"fmt"
	"github.com/autom8ter/thermomatic/internal/common"
)

// Decode returns the IMEI code contained in the first 15 bytes of b.
//
// In case b isn't strictly composed of digits, the returned error will be
// ErrInvalid.
//
// In case b's checksum is wrong, the returned error will be ErrChecksum.
//
// Decode does NOT allocate under any condition. Additionally, it panics if b
// isn't at least 15 bytes long.
func Decode(b []byte) (uint64, error) {
	var (
		sum    = uint64(0)
		actual = uint64(0)
		double = false
	)

	if len(b) < common.MinImeiLength {
		//code must be at least 15 bytes long.
		panic(common.ErrInvalidImei)
	}

	//iterate over each digit in the byteslice
	for i := 0; i < common.MinImeiLength; i++ {
		digit := uint64(b[i] - common.ASCIIZero)
		//each digit should be between 0-9
		if digit > 9 {
			digit = digit - 9
			return 0, common.Wrap(common.ErrInvalidImei, fmt.Sprintf("invalid digit: %d", digit))
		}
		//base10
		actual = (uint64(10) * actual) + digit

		//skip last digit when calculating sum for luhn validation ref: https://en.wikipedia.org/wiki/International_Mobile_Equipment_Identity.
		if i == 14 {
			continue
		}
		if double {
			digit = digit * 2
		}
		if digit >= 10 {
			digit = digit - 9
		}
		sum += digit
		//double every other digit
		double = !double
	}
	//validate using Luhn algorithm
	if ((10 - (sum % 10)) % 10) != uint64(b[14]-common.ASCIIZero) {
		return 0, common.Wrap(common.ErrChecksum, fmt.Sprintf("payload = %s sum = %v", string(b), sum))
	}
	return actual, nil
}
