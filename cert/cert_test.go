package cert

import (
	"comradequinn/hflow/log"
	"crypto/tls"
	"strconv"
	"testing"
)

// noOptimize is used to assign unneeded results to so as to
// prevent the compiler attempting to optimising the program
// by removing a seemingly redundant function calls
var noOptimize any

func Benchmark_For(b *testing.B) {
	noOptimize, _ = For(&tls.ClientHelloInfo{ServerName: "warm-up.func"})
	log.SetVerbosity(0)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		noOptimize, _ = For(&tls.ClientHelloInfo{ServerName: "domain" + strconv.Itoa(b.N%10) + ".com"})
	}
}

func TestCertGet(t *testing.T) {
	cert, err := For(&tls.ClientHelloInfo{ServerName: "domain.com"})

	if err != nil {
		t.Fatalf("expected no error after cert generation, got: [%v]", err)
	}

	if cert == nil {
		t.Fatalf("expected certificate after cert generation, got: nil")
	}
}
