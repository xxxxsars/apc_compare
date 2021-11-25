package main

import (
	"apc_compare/pkg/compare"
	"testing"
)

// go test -v -bench=. -run=none .
func BenchmarkSpcCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		compare.SpaceCompare()
	}
}
