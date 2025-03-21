package kretry

import (
	"fmt"
)

func ExampleBackoff() {
	b := NewBackoff()
	for i := 0; i < 10; i++ {
		fmt.Println(b.Duration())
	}
	// Output:
	// 100ms
	// 200ms
	// 400ms
	// 800ms
	// 1.6s
	// 3.2s
	// 6.4s
	// 10s
	// 10s
	// 10s

}

func ExampleBackoff_WithFactor() {
	b := NewBackoff(WithFactor(4))
	for i := 0; i < 10; i++ {
		fmt.Println(b.Duration())
	}
	// Output:
	// 100ms
	// 400ms
	// 1.6s
	// 6.4s
	// 10s
	// 10s
	// 10s
	// 10s
	// 10s
	// 10s
}

func ExampleBackoff_WithJitter() {
	b := NewBackoff(WithJitter(true))
	for i := 0; i < 10; i++ {
		fmt.Println(b.Duration())
	}
	// Output:
	// 可能的输出
	// 100ms
	// 130.324613ms
	// 290.318078ms
	// 282.896353ms
	// 1.290019019s
	// 1.237657307s
	// 3.30529213s
	// 1.420308306s
	// 10s
	// 10s
}
