package num

import "testing"

func TestGenericNumberAggregates(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "signed integers",
			run: func(t *testing.T) {
				if got := SumNumber[int](-2, 5, 7); got != 10 {
					t.Fatalf("SumNumber[int] = %v", got)
				}
				if got := AvgNumber[int](-2, 5, 7); got != 10.0/3.0 {
					t.Fatalf("AvgNumber[int] = %v", got)
				}
				if got := MinInteger[int](-2, 5); got != -2 {
					t.Fatalf("MinInteger[int] = %d", got)
				}
				if got := MinIntegers[int](4, -8, 2); got != -8 {
					t.Fatalf("MinIntegers[int] = %d", got)
				}
				if got := MaxInteger[int](-2, 5); got != 5 {
					t.Fatalf("MaxInteger[int] = %d", got)
				}
				if got := MaxIntegers[int](4, -8, 2); got != 4 {
					t.Fatalf("MaxIntegers[int] = %d", got)
				}
			},
		},
		{
			name: "unsigned integers",
			run: func(t *testing.T) {
				if got := SumNumber[uint](2, 5, 7); got != 14 {
					t.Fatalf("SumNumber[uint] = %v", got)
				}
				if got := AvgNumber[uint](2, 5, 8); got != 5 {
					t.Fatalf("AvgNumber[uint] = %v", got)
				}
				if got := MinIntegers[uint](4, 8, 2); got != 2 {
					t.Fatalf("MinIntegers[uint] = %d", got)
				}
				if got := MaxIntegers[uint](4, 8, 2); got != 8 {
					t.Fatalf("MaxIntegers[uint] = %d", got)
				}
			},
		},
		{
			name: "floats",
			run: func(t *testing.T) {
				if got := SumNumber[float64](1.25, 2.5, -0.75); got != 3 {
					t.Fatalf("SumNumber[float64] = %v", got)
				}
				if got := AvgNumber[float32](1.5, 2.5); got != 2 {
					t.Fatalf("AvgNumber[float32] = %v", got)
				}
				if got := MinFloat64(1.25, -3.5); got != -3.5 {
					t.Fatalf("MinFloat64 = %v", got)
				}
				if got := MaxFloat64(1.25, -3.5); got != 1.25 {
					t.Fatalf("MaxFloat64 = %v", got)
				}
				if got := MinFloat64s(3.5, -1.25, 2); got != -1.25 {
					t.Fatalf("MinFloat64s = %v", got)
				}
				if got := MaxFloat64s(3.5, -1.25, 2); got != 3.5 {
					t.Fatalf("MaxFloat64s = %v", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func TestGenericNumberEmptyInputsReturnZero(t *testing.T) {
	if got := AvgNumber[int](); got != 0 {
		t.Fatalf("AvgNumber empty = %v", got)
	}
	if got := MinIntegers[int](); got != 0 {
		t.Fatalf("MinIntegers empty = %d", got)
	}
	if got := MaxIntegers[int](); got != 0 {
		t.Fatalf("MaxIntegers empty = %d", got)
	}
	if got := MinFloat64s(); got != 0 {
		t.Fatalf("MinFloat64s empty = %v", got)
	}
	if got := MaxFloat64s(); got != 0 {
		t.Fatalf("MaxFloat64s empty = %v", got)
	}
}
