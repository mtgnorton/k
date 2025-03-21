package kmath

import "testing"

func TestMax(t *testing.T) {
	if Max(1, 2) != 2 {
		t.Error("Max(1, 2) != 2")
	}
	if Max(1.1, 2.2) != 2.2 {
		t.Error("Max(1.1, 2.2) != 2.2")
	}
}

func TestMin(t *testing.T) {
	if Min(1, 2) != 1 {
		t.Error("Min(1, 2) != 1")
	}
	if Min(1.1, 2.2) != 1.1 {
		t.Error("Min(1.1, 2.2) != 1.1")
	}
}

func TestRound(t *testing.T) {
	if Round(1.234, 2) != 1.23 {
		t.Error("Round(1.234, 2) != 1.23")
	}
}

func TestFloor(t *testing.T) {
	if Floor(1.6) != 1 {
		t.Error("Floor(1.6) != 1")
	}
}

func TestCeil(t *testing.T) {
	if Ceil(1.2, 0) != 2 {
		t.Error("Ceil(1.2, 0) != 2")
	}
}

func TestAbs(t *testing.T) {
	if Abs(-1) != 1 {
		t.Error("Abs(-1) != 1")
	}
	if Abs(-1.1) != 1.1 {
		t.Error("Abs(-1.1) != 1.1")
	}
}

func TestPow(t *testing.T) {
	if Pow(2, 3) != 8 {
		t.Error("Pow(2, 3) != 8")
	}
}

func TestSqrt(t *testing.T) {
	if Sqrt(4) != 2 {
		t.Error("Sqrt(4) != 2")
	}
}

func TestRandInt(t *testing.T) {
	min, max := 1, 10
	for i := 0; i < 100; i++ {
		n := RandInt(min, max)
		if n < min || n > max {
			t.Errorf("RandInt(%d, %d) = %d", min, max, n)
		}
	}
}

func TestRandFloat(t *testing.T) {
	min, max := 1.0, 10.0
	for i := 0; i < 100; i++ {
		n := RandFloat(min, max)
		if n < min || n > max {
			t.Errorf("RandFloat(%f, %f) = %f", min, max, n)
		}
	}
}
