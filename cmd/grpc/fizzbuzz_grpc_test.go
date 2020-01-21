package main

import (
	"testing"
	"os"
	"flag"
)

func TestMain(m *testing.M) {
	flag.Parse()
    os.Exit(m.Run())
}