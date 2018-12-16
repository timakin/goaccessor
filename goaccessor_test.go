package goaccessor_test

import (
	"log"
	"testing"

	"github.com/timakin/goaccessor"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	log.Printf("%+v", analysistest.Run(t, testdata, goaccessor.Analyzer, "testpackage")[0])
}
