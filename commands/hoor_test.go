package commands

import (
	"testing"
)

func TestPersianizeNumbers(t *testing.T) {
	type data struct {
		TestName      string
		NormalString  string
		PersianString string
	}

	tests := []data{
		{"Alphanumeric", "سلام جهان 12", "سلام جهان ۱۲"},
		{"Numeric", "0123456789", "۰۱۲۳۴۵۶۷۸۹"},
		{"RepititiveAlphanumeric", "11 foo 153 بار", "۱۱ foo ۱۵۳ بار"},
	}

	for i, test := range tests {
		output := persianizeNumbers(test.NormalString)
		if output != test.PersianString {
			t.Errorf("Test #%d %s: expected %q, got %q", i, test.TestName, test.PersianString, output)
		}
	}
}
