package hw10programoptimization

import (
	"bytes"
	"testing"
)

func BenchmarkGetDomainStat(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = GetDomainStat(bytes.NewBufferString(data), "biz")
	}
}
