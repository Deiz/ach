// Licensed to The Moov Authors under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. The Moov Authors licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package ach

import (
	"fmt"
	"testing"
)

func TestValidators__checkDigit(t *testing.T) {
	cases := map[string]int{
		// invalid
		"":          -1,
		"123456":    -1,
		"1a8ab":     -1,
		"0730002a":  -1,
		"0730A002":  -1,
		"YYYYYYYYY": -1, // users often mask ABA numbers
		// valid
		"07300022": 8, // Wells Fargo - Iowa
		"10200007": 6, // Wells Fargo - Colorado
	}

	v := validator{}
	for rtn, check := range cases {
		answer := v.CalculateCheckDigit(rtn)
		if check != answer {
			t.Errorf("input=%s answer=%d expected=%d", rtn, answer, check)
		}
		if err := CheckRoutingNumber(fmt.Sprintf("%s%d", rtn, check)); err != nil && check >= 0 {
			t.Errorf("input=%s answer=%d expected=%d: %v", rtn, answer, check, err)
		}
	}
}

func TestValidators__isCreditCardYear(t *testing.T) {
	cases := map[string]bool{
		// invalid (or out of range)
		"10": false,
		"00": false,
		"51": false,
		"17": false,
		// valid
		"20": true,
		"19": true,
	}
	v := validator{}
	for yy, valid := range cases {
		err := v.isCreditCardYear(yy)
		if valid && err != nil {
			t.Errorf("yy=%s failed: %v", yy, err)
		}
		if !valid && err == nil {
			t.Errorf("yy=%s should have failed", yy)
		}
	}
}

func TestValidators__validateSimpleDate(t *testing.T) {
	cases := map[string]string{
		// invalid
		"":       "",
		"01":     "",
		"001520": "", // no 15th month
		"001240": "", // no 40th Day
		"190001": "", // no 0th month
		"190100": "", // no 0th day
		// valid
		"190101": "190101", // Jan 1st
		"201231": "201231", // Dec 31st
		"220731": "220731", // July 31st
		"350430": "350430", // April 30th
		"500229": "500229", // Feb 29th
	}

	v := validator{}
	for input, expected := range cases {
		answer := v.validateSimpleDate(input)
		if expected != answer {
			t.Errorf("input=%q got=%q expected=%q", input, answer, expected)
		}
	}
}

func TestValidators__validateSimpleTime(t *testing.T) {
	cases := map[string]string{
		// invalid
		"":       "",
		"01":     "",
		"012":    "",
		"123142": "",
		// valid
		"0000": "0000",
		"0100": "0100",
		"2359": "2359",
		"1201": "1201",
		"1238": "1238",
	}
	v := validator{}
	for input, expected := range cases {
		answer := v.validateSimpleTime(input)
		if expected != answer {
			t.Errorf("input=%q got=%q expected=%q", input, answer, expected)
		}
	}
}

// Ensure that ASCII characters from 0x20 to 0x7E are considered alphanumeric.
func TestValidators__isAlphanumeric(t *testing.T) {
	v := validator{}

	for i := 0; i < 255; i++ {
		err := v.isAlphanumeric(string(rune(i)))

		shouldError := i < 0x20 || i > 0x7E

		if shouldError && err == nil {
			t.Errorf("expected rune %x to be non-alphanumeric", i)
		} else if !shouldError && err != nil {
			t.Errorf("expected rune %x to be alphanumeric: %v", i, err)
		}
	}
}

// Ensure that ASCII characters from 0x20 to 0x60 and 0x7B to 0x7E are considered upper case alphanumeric.
func TestValidators__isUpperAlphanumeric(t *testing.T) {
	v := validator{}

	for i := 0; i < 255; i++ {
		chr := string(rune(i))
		err := v.isUpperAlphanumeric(chr)

		shouldError := i < 0x20 || i > 0x7E || (i > 0x60 && i < 0x7B)

		if shouldError && err == nil {
			t.Errorf("expected rune %x (%s) to be non-alphanumeric", i, chr)
		} else if !shouldError && err != nil {
			t.Errorf("expected rune %x (%s) to be alphanumeric: %v", i, chr, err)
		}
	}
}
