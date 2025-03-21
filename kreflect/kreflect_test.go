package kreflect

import "testing"

func TestIsNil(t *testing.T) {
	if !IsNil(nil) {
		t.Error("IsNil(nil) != true")
	}
	if IsNil(1) {
		t.Error("IsNil(1) != false")
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		input    any
		expected string
	}{
		{1, "1"},
		{1.1, "1.1"},
		{true, "true"},
		{"hello", "hello"},
		{[]byte("hello"), "hello"},
		{nil, ""},
	}

	for _, test := range tests {
		if got := ToString(test.input); got != test.expected {
			t.Errorf("ToString(%v) = %v, want %v", test.input, got, test.expected)
		}
	}
}
