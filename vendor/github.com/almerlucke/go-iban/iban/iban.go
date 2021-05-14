package iban

import (
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

// CountrySettings contains length for IBAN and format for BBAN
type CountrySettings struct {
	// Length of IBAN code for this country
	Length int

	// Format of BBAN part of IBAN for this country
	Format string
}

// IBAN struct
type IBAN struct {
	// Full code
	Code string

	// Full code prettyfied for printing on paper
	PrintCode string

	// Country code
	CountryCode string

	// Check digits
	CheckDigits string

	// Country settings
	CountrySettings *CountrySettings

	// Country specific bban part
	BBAN string
}

/*
	Taken from http://www.tbg5-finance.org/ code example
*/
var countries = map[string]CountrySettings{
	"AD": CountrySettings{Length: 24, Format: "F04F04A12"},
	"AE": CountrySettings{Length: 23, Format: "F03F16"},
	"AL": CountrySettings{Length: 28, Format: "F08A16"},
	"AT": CountrySettings{Length: 20, Format: "F05F11"},
	"AZ": CountrySettings{Length: 28, Format: "U04A20"},
	"BA": CountrySettings{Length: 20, Format: "F03F03F08F02"},
	"BE": CountrySettings{Length: 16, Format: "F03F07F02"},
	"BG": CountrySettings{Length: 22, Format: "U04F04F02A08"},
	"BH": CountrySettings{Length: 22, Format: "U04A14"},
	"BR": CountrySettings{Length: 29, Format: "F08F05F10U01A01"},
	"CH": CountrySettings{Length: 21, Format: "F05A12"},
	"CR": CountrySettings{Length: 21, Format: "F03F14"},
	"CY": CountrySettings{Length: 28, Format: "F03F05A16"},
	"CZ": CountrySettings{Length: 24, Format: "F04F06F10"},
	"DE": CountrySettings{Length: 22, Format: "F08F10"},
	"DK": CountrySettings{Length: 18, Format: "F04F09F01"},
	"DO": CountrySettings{Length: 28, Format: "U04F20"},
	"EE": CountrySettings{Length: 20, Format: "F02F02F11F01"},
	"ES": CountrySettings{Length: 24, Format: "F04F04F01F01F10"},
	"FI": CountrySettings{Length: 18, Format: "F06F07F01"},
	"FO": CountrySettings{Length: 18, Format: "F04F09F01"},
	"FR": CountrySettings{Length: 27, Format: "F05F05A11F02"},
	"GB": CountrySettings{Length: 22, Format: "U04F06F08"},
	"GE": CountrySettings{Length: 22, Format: "U02F16"},
	"GI": CountrySettings{Length: 23, Format: "U04A15"},
	"GL": CountrySettings{Length: 18, Format: "F04F09F01"},
	"GR": CountrySettings{Length: 27, Format: "F03F04A16"},
	"GT": CountrySettings{Length: 28, Format: "A04A20"},
	"HR": CountrySettings{Length: 21, Format: "F07F10"},
	"HU": CountrySettings{Length: 28, Format: "F03F04F01F15F01"},
	"IE": CountrySettings{Length: 22, Format: "U04F06F08"},
	"IL": CountrySettings{Length: 23, Format: "F03F03F13"},
	"IS": CountrySettings{Length: 26, Format: "F04F02F06F10"},
	"IT": CountrySettings{Length: 27, Format: "U01F05F05A12"},
	"JO": CountrySettings{Length: 30, Format: "U04F04A18"},
	"KW": CountrySettings{Length: 30, Format: "U04A22"},
	"KZ": CountrySettings{Length: 20, Format: "F03A13"},
	"LB": CountrySettings{Length: 28, Format: "F04A20"},
	"LC": CountrySettings{Length: 32, Format: "U04A24"},
	"LI": CountrySettings{Length: 21, Format: "F05A12"},
	"LT": CountrySettings{Length: 20, Format: "F05F11"},
	"LU": CountrySettings{Length: 20, Format: "F03A13"},
	"LV": CountrySettings{Length: 21, Format: "U04A13"},
	"MC": CountrySettings{Length: 27, Format: "F05F05A11F02"},
	"MD": CountrySettings{Length: 24, Format: "A20"},
	"ME": CountrySettings{Length: 22, Format: "F03F13F02"},
	"MK": CountrySettings{Length: 19, Format: "F03A10F02"},
	"MR": CountrySettings{Length: 27, Format: "F05F05F11F02"},
	"MT": CountrySettings{Length: 31, Format: "U04F05A18"},
	"MU": CountrySettings{Length: 30, Format: "U04F02F02F12F03U03"},
	"NL": CountrySettings{Length: 18, Format: "U04F10"},
	"NO": CountrySettings{Length: 15, Format: "F04F06F01"},
	"PK": CountrySettings{Length: 24, Format: "U04A16"},
	"PL": CountrySettings{Length: 28, Format: "F08F16"},
	"PS": CountrySettings{Length: 29, Format: "U04A21"},
	"PT": CountrySettings{Length: 25, Format: "F04F04F11F02"},
	"QA": CountrySettings{Length: 29, Format: "U04A21"},
	"RO": CountrySettings{Length: 24, Format: "U04A16"},
	"RS": CountrySettings{Length: 22, Format: "F03F13F02"},
	"SA": CountrySettings{Length: 24, Format: "F02A18"},
	"SC": CountrySettings{Length: 31, Format: "U04F02F02F16U03"},
	"SE": CountrySettings{Length: 24, Format: "F03F16F01"},
	"SI": CountrySettings{Length: 19, Format: "F05F08F02"},
	"SK": CountrySettings{Length: 24, Format: "F04F06F10"},
	"SM": CountrySettings{Length: 27, Format: "U01F05F05A12"},
	"ST": CountrySettings{Length: 25, Format: "F08F11F02"},
	"TL": CountrySettings{Length: 23, Format: "F03F14F02"},
	"TN": CountrySettings{Length: 24, Format: "F02F03F13F02"},
	"TR": CountrySettings{Length: 26, Format: "F05A01A16"},
	"UA": CountrySettings{Length: 29, Format: "F06A19"},
	"VG": CountrySettings{Length: 24, Format: "U04F16"},
	"XK": CountrySettings{Length: 20, Format: "F04F10F02"},
}

func validateCheckDigits(iban string) error {
	// Move the four initial characters to the end of the string
	iban = iban[4:] + iban[:4]

	// Replace each letter in the string with two digits, thereby expanding the string, where A = 10, B = 11, ..., Z = 35
	mods := ""
	for _, c := range iban {
		// Get character code point value
		i := int(c)

		// Check if c is characters A-Z (codepoint 65 - 90)
		if i > 64 && i < 91 {
			// A=10, B=11 etc...
			i -= 55
			// Add int as string to mod string
			mods += strconv.Itoa(i)
		} else {
			mods += string(c)
		}
	}

	// Create bignum from mod string and perform module
	bigVal, success := new(big.Int).SetString(mods, 10)
	if !success {
		return errors.New("IBAN check digits validation failed")
	}

	modVal := new(big.Int).SetInt64(97)
	resVal := new(big.Int).Mod(bigVal, modVal)

	// Check if module is equal to 1
	if resVal.Int64() != 1 {
		return errors.New("IBAN has incorrect check digits")
	}

	return nil
}

func validateBasicBankAccountNumber(bban string, format string) error {
	// Format regex to get parts
	frx, err := regexp.Compile(`[ABCFLUW]\d{2}`)
	if err != nil {
		return fmt.Errorf("Failed to validate BBAN: %v", err.Error())
	}

	// Get format part strings
	fps := frx.FindAllString(format, -1)

	// Create regex from format parts
	bbr := ""

	for _, ps := range fps {
		switch ps[:1] {
		case "F":
			bbr += "[0-9]"
		case "L":
			bbr += "[a-z]"
		case "U":
			bbr += "[A-Z]"
		case "A":
			bbr += "[0-9A-Za-z]"
		case "B":
			bbr += "[0-9A-Z]"
		case "C":
			bbr += "[A-Za-z]"
		case "W":
			bbr += "[0-9a-z]"
		}

		// Get repeat factor for group
		repeat, atoiErr := strconv.Atoi(ps[1:])
		if atoiErr != nil {
			return fmt.Errorf("Failed to validate BBAN: %v", atoiErr.Error())
		}

		// Add to regex
		bbr += fmt.Sprintf("{%d}", repeat)
	}

	// Compile regex and validate BBAN
	bbrx, err := regexp.Compile(bbr)
	if err != nil {
		return fmt.Errorf("Failed to validate BBAN: %v", err.Error())
	}

	if !bbrx.MatchString(bban) {
		return errors.New("BBAN part of IBAN is not formatted according to country specification")
	}

	return nil
}

// NewIBAN create new IBAN with validation
func NewIBAN(s string) (*IBAN, error) {
	iban := IBAN{}

	// Prepare string: remove spaces and convert to upper case
	s = strings.ToUpper(strings.Replace(s, " ", "", -1))
	iban.Code = s

	// Validate characters
	r, err := regexp.Compile(`^[0-9A-Z]*$`)
	if err != nil {
		return nil, fmt.Errorf("Failed to validate IBAN: %v", err.Error())
	}

	if !r.MatchString(s) {
		return nil, errors.New("IBAN can contain only alphanumeric characters")
	}

	// Get country code and check digits
	r, err = regexp.Compile(`^\D\D\d\d`)
	if err != nil {
		return nil, fmt.Errorf("Failed to validate IBAN: %v", err.Error())
	}

	hs := r.FindString(s)
	if hs == "" {
		return nil, errors.New("IBAN must start with country code (2 characters) and check digits (2 digits)")
	}

	iban.CountryCode = hs[0:2]
	iban.CheckDigits = hs[2:4]

	// Get country settings for country code
	cs, ok := countries[iban.CountryCode]
	if !ok {
		return nil, fmt.Errorf("Unsupported country code %v", iban.CountryCode)
	}

	iban.CountrySettings = &cs

	// Validate code length
	if len(s) != cs.Length {
		return nil, fmt.Errorf("IBAN length %d does not match length %d specified for country code %v", len(s), cs.Length, iban.CountryCode)
	}

	// Set and validate BBAN part, the part after the language code and check digits
	iban.BBAN = s[4:]

	err = validateBasicBankAccountNumber(iban.BBAN, iban.CountrySettings.Format)
	if err != nil {
		return nil, err
	}

	// Validate check digits with mod97
	err = validateCheckDigits(iban.Code)
	if err != nil {
		return nil, err
	}

	// Generate print code from code (splits code in sections of 4 characters)
	prc := ""
	for len(s) > 4 {
		prc += s[:4] + " "
		s = s[4:]
	}

	iban.PrintCode = prc + s

	return &iban, nil
}
