package main

import (
	"github.com/oribe1115/gentest"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(gentest.Analyzer) }

