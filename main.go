package main

import (
	"github.com/lazyjean/sla2/app"
	"github.com/lazyjean/sla2/internal/wire"
)

func main() {
	app.Run(wire.InitializeApp())
}
