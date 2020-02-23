package brainwave_yagpdb

import (
	"testing"

	"github.com/jonas747/dcmd"
)

var _ dcmd.ArgType = (*ActionArg)(nil)

func TestMatches(t *testing.T) {
	tests := []struct {
		input string
		ok    bool
	}{
		{"lead", true},
		{"Lead", true},
		{"LEAD", true},
		{"", false},
		{"asdf", false},
	}

	aa := ActionArg{}
	for _, test := range tests {
		ok := aa.Matches(nil, test.input)
		if ok != test.ok {
			t.Errorf("%v; got %t", test, ok)
		}
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		input  string
		action Action
		ok     bool
	}{
		{"lead", Lead, true},
		{"Lead", Lead, true},
		{"LEAD", Lead, true},
		{"", Invalid, false},
		{"asdf", Invalid, false},
	}

	aa := ActionArg{}
	for _, test := range tests {
		a, err := aa.Parse(nil, test.input, nil)
		if a != test.action || ((err == nil) != test.ok) {
			t.Errorf("%v; got %s %t", test, a, err == nil)
		}
	}
}
