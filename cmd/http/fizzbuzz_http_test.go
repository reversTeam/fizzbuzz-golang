package main

import (
	"testing"
	"os"
	"flag"
)

func TestMain(m *testing.M) {
	flag.Parse()
	defer os.Exit(m.Run())
}
