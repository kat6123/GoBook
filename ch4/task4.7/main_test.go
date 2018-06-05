package main

import "testing"

func BenchmarkReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reverse([]byte("Привет"))
	}
}

func BenchmarkReverseV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reverseV2([]byte("Привет"))
	}
}
