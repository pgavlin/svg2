package main

import (
	"encoding/xml"
	"log"
	"os"

	"github.com/pgavlin/svg"
)

func main() {
	var doc svg.SVG
	if err := xml.NewDecoder(os.Stdin).Decode(&doc); err != nil {
		log.Fatal(err)
	}

	ctx := svg.NewContext(&doc)
	if err := svg.Render(ctx, &doc); err != nil {
		log.Fatal(err)
	}

	if err := ctx.EncodePNG(os.Stdout); err != nil {
		log.Fatal(err)
	}
}
