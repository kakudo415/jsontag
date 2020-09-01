package jsontag_test

import (
	"testing"

	"github.com/kakudo415/jsontag"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, jsontag.Analyzer, "a")
}

