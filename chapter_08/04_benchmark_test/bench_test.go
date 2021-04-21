package main

import (
	"testing"
)

// run becnhmark tests: go test -v -cover -short -bench .
// ignore any tests not maching 'x': go test -run x -bench .
func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		decode("post.json")
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		unmarshal("post.json")
	}
}
