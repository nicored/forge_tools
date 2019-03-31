package properties

import (
	"fmt"
	"testing"
	"time"
)

func TestProperties_Run(t *testing.T) {
	tStart := time.Now()
	p := NewProperties("./samples/big_sample")
	p.Run()
	tEnd := time.Now()
	fmt.Println(tEnd.Sub(tStart).String())
}

func BenchmarkProperties_Run(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		p := NewProperties("./samples/big_sample")
		p.Run()
	}
}
