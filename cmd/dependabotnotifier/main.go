package main

import (
	"context"

	"github.com/navikt/dependabotnotifier/internal/cmd/dependabotnotifier"
)

func main() {
	dependabotnotifier.Run(context.Background())
}
