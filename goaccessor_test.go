package goaccessor_test

import (
	"testing"

	"github.com/timakin/goaccessor"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, goaccessor.Analyzer, "testpackage")
}
