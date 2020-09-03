package main

import (
	"github.com/kakudo415/jsontag"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(jsontag.Analyzer) }
