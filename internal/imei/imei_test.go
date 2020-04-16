package imei_test

import (
	"github.com/autom8ter/thermomatic/internal/imei"
	"testing"
)

//TestDecode fails if the output of Decode doesnt match the Expect value. It also fails if there are any allocations.
func TestDecode(t *testing.T) {
	tests := []struct {
		Name   string
		Code   []byte
		Expect uint64
	}{

		{
			Name:   "imei code (1)",
			Code:   []byte("450711608247968"),
			Expect: 450711608247968,
		},
		{
			Name:   "longer, but valid imei code (2)",
			Code:   []byte("529573786277564560"),
			Expect: 529573786277564,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			apr := testing.AllocsPerRun(5000, func() {
				actual, err := imei.Decode(test.Code)
				if err != nil {
					t.Fatalf("unexpected error = %s\n", err.Error())
				}
				if test.Expect != actual {
					t.Fatalf("actual: %v expected: %v\n", test.Code, test.Expect)
				}
			})
			if apr > 0 {
				t.Fatal("allocations per run is greater than zero!")
			}
		})
	}
}

//go test -v -bench=.
//BenchmarkDecode-12      50000000                23.9 ns/op             0 B/op          0 allocs/op
func BenchmarkDecode(b *testing.B) {
	b.ReportAllocs()
	code := []byte("450711608247968")
	for i := 0; i < b.N; i++ {
		_, err := imei.Decode(code)
		if err != nil {
			b.Fatalf("unexpected error = %s\n", err)
		}
	}
}
