package myfx

import (
	"os"
	"strconv"
)

type Verbose interface {
	Set(verbose bool)
	Get() bool
}

type verbose struct {
	verbose bool
}

func (v *verbose) Set(verbose bool) {
	v.verbose = verbose
}

func (v *verbose) Get() bool {
	return v.verbose
}

func NewVerbose() Verbose {
	vb := &verbose{}
	vb.verbose, _ = strconv.ParseBool(os.Getenv("VERBOSE"))
	return vb
}
