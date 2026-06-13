package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

const systemName = "xquare-deploy"

func main() {
	if err := run(os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func run(output io.Writer) error {
	if _, err := fmt.Fprintln(output, systemName); err != nil {
		return fmt.Errorf("write system name: %w", err)
	}

	return nil
}
