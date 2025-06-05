package main

import (
	"github.com/lazyjean/sla2/internal/wire"
)

func main() {
	if err := wire.Run(); err != nil {
		panic(err)
	}
}
