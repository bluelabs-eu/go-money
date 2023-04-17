package money

import (
	"testing"
)

func TestFormatter_Format(t *testing.T) {
	tcs := []struct {
		fraction int
		decimal  string
		thousand string
		grapheme string
		template string
		amount   int64
		expected string
	}{
		{2, ".", ",", "$", "1 $", 0, "0.00 $"},
		{2, ".", ",", "$", "1 $", 1, "0.01 $"},
		{2, ".", ",", "$", "1 $", 12, "0.12 $"},
		{2, ".", ",", "$", "1 $", 123, "1.23 $"},
		{2, ".", ",", "$", "1 $", 1234, "12.34 $"},
		{2, ".", ",", "$", "1 $", 12345, "123.45 $"},
		{2, ".", ",", "$", "1 $", 123456, "1,234.56 $"},
		{2, ".", ",", "$", "1 $", 1234567, "12,345.67 $"},
		{2, ".", ",", "$", "1 $", 12345678, "123,456.78 $"},
		{2, ".", ",", "$", "1 $", 123456789, "1,234,567.89 $"},

		{2, ".", ",", "$", "1 $", -1, "-0.01 $"},
		{2, ".", ",", "$", "1 $", -12, "-0.12 $"},
		{2, ".", ",", "$", "1 $", -123, "-1.23 $"},
		{2, ".", ",", "$", "1 $", -1234, "-12.34 $"},
		{2, ".", ",", "$", "1 $", -12345, "-123.45 $"},
		{2, ".", ",", "$", "1 $", -123456, "-1,234.56 $"},
		{2, ".", ",", "$", "1 $", -1234567, "-12,345.67 $"},
		{2, ".", ",", "$", "1 $", -12345678, "-123,456.78 $"},
		{2, ".", ",", "$", "1 $", -123456789, "-1,234,567.89 $"},

		{3, ".", "", "$", "1 $", 1, "0.001 $"},
		{3, ".", "", "$", "1 $", 12, "0.012 $"},
		{3, ".", "", "$", "1 $", 123, "0.123 $"},
		{3, ".", "", "$", "1 $", 1234, "1.234 $"},
		{3, ".", "", "$", "1 $", 12345, "12.345 $"},
		{3, ".", "", "$", "1 $", 123456, "123.456 $"},
		{3, ".", "", "$", "1 $", 1234567, "1234.567 $"},
		{3, ".", "", "$", "1 $", 12345678, "12345.678 $"},
		{3, ".", "", "$", "1 $", 123456789, "123456.789 $"},

		{2, ".", ",", "£", "$1", 1, "£0.01"},
		{2, ".", ",", "£", "$1", 12, "£0.12"},
		{2, ".", ",", "£", "$1", 123, "£1.23"},
		{2, ".", ",", "£", "$1", 1234, "£12.34"},
		{2, ".", ",", "£", "$1", 12345, "£123.45"},
		{2, ".", ",", "£", "$1", 123456, "£1,234.56"},
		{2, ".", ",", "£", "$1", 1234567, "£12,345.67"},
		{2, ".", ",", "£", "$1", 12345678, "£123,456.78"},
		{2, ".", ",", "£", "$1", 123456789, "£1,234,567.89"},

		{0, ".", ",", "NT$", "$1", 1, "NT$1"},
		{0, ".", ",", "NT$", "$1", 12, "NT$12"},
		{0, ".", ",", "NT$", "$1", 123, "NT$123"},
		{0, ".", ",", "NT$", "$1", 1234, "NT$1,234"},
		{0, ".", ",", "NT$", "$1", 12345, "NT$12,345"},
		{0, ".", ",", "NT$", "$1", 123456, "NT$123,456"},
		{0, ".", ",", "NT$", "$1", 1234567, "NT$1,234,567"},
		{0, ".", ",", "NT$", "$1", 12345678, "NT$12,345,678"},
		{0, ".", ",", "NT$", "$1", 123456789, "NT$123,456,789"},

		{0, ".", ",", "NT$", "$1", -1, "-NT$1"},
		{0, ".", ",", "NT$", "$1", -12, "-NT$12"},
		{0, ".", ",", "NT$", "$1", -123, "-NT$123"},
		{0, ".", ",", "NT$", "$1", -1234, "-NT$1,234"},
		{0, ".", ",", "NT$", "$1", -12345, "-NT$12,345"},
		{0, ".", ",", "NT$", "$1", -123456, "-NT$123,456"},
		{0, ".", ",", "NT$", "$1", -1234567, "-NT$1,234,567"},
		{0, ".", ",", "NT$", "$1", -12345678, "-NT$12,345,678"},
		{0, ".", ",", "NT$", "$1", -123456789, "-NT$123,456,789"},
	}

	for _, tc := range tcs {
		formatter := NewFormatter(tc.fraction, tc.decimal, tc.thousand, tc.grapheme, tc.template)
		r := formatter.Format(tc.amount)

		if r != tc.expected {
			t.Errorf("Expected %d formatted to be %s got %s", tc.amount, tc.expected, r)
		}
	}
}

func TestFormatter_FormatAmount(t *testing.T) {
	tcs := []struct {
		fraction int
		decimal  string
		thousand string
		grapheme string
		template string
		amount   int64
		expected string
	}{
		{2, ".", ",", "$", "1 $", 0, "0.00"},
		{2, ".", ",", "$", "1 $", 1, "0.01"},
		{2, ".", ",", "$", "1 $", 12, "0.12"},
		{2, ".", ",", "$", "1 $", 123, "1.23"},
		{2, ".", ",", "$", "1 $", 1234, "12.34"},
		{2, ".", ",", "$", "1 $", 12345, "123.45"},
		{2, ".", ",", "$", "1 $", 123456, "1,234.56"},
		{2, ".", ",", "$", "1 $", 1234567, "12,345.67"},
		{2, ".", ",", "$", "1 $", 12345678, "123,456.78"},
		{2, ".", ",", "$", "1 $", 123456789, "1,234,567.89"},

		{2, ".", ",", "$", "1 $", -1, "-0.01"},
		{2, ".", ",", "$", "1 $", -12, "-0.12"},
		{2, ".", ",", "$", "1 $", -123, "-1.23"},
		{2, ".", ",", "$", "1 $", -1234, "-12.34"},
		{2, ".", ",", "$", "1 $", -12345, "-123.45"},
		{2, ".", ",", "$", "1 $", -123456, "-1,234.56"},
		{2, ".", ",", "$", "1 $", -1234567, "-12,345.67"},
		{2, ".", ",", "$", "1 $", -12345678, "-123,456.78"},
		{2, ".", ",", "$", "1 $", -123456789, "-1,234,567.89"},

		{3, ".", "", "$", "1 $", 1, "0.001"},
		{3, ".", "", "$", "1 $", 12, "0.012"},
		{3, ".", "", "$", "1 $", 123, "0.123"},
		{3, ".", "", "$", "1 $", 1234, "1.234"},
		{3, ".", "", "$", "1 $", 12345, "12.345"},
		{3, ".", "", "$", "1 $", 123456, "123.456"},
		{3, ".", "", "$", "1 $", 1234567, "1234.567"},
		{3, ".", "", "$", "1 $", 12345678, "12345.678"},
		{3, ".", "", "$", "1 $", 123456789, "123456.789"},

		{2, ".", ",", "£", "$1", 1, "0.01"},
		{2, ".", ",", "£", "$1", 12, "0.12"},
		{2, ".", ",", "£", "$1", 123, "1.23"},
		{2, ".", ",", "£", "$1", 1234, "12.34"},
		{2, ".", ",", "£", "$1", 12345, "123.45"},
		{2, ".", ",", "£", "$1", 123456, "1,234.56"},
		{2, ".", ",", "£", "$1", 1234567, "12,345.67"},
		{2, ".", ",", "£", "$1", 12345678, "123,456.78"},
		{2, ".", ",", "£", "$1", 123456789, "1,234,567.89"},

		{0, ".", ",", "NT$", "$1", 1, "1"},
		{0, ".", ",", "NT$", "$1", 12, "12"},
		{0, ".", ",", "NT$", "$1", 123, "123"},
		{0, ".", ",", "NT$", "$1", 1234, "1,234"},
		{0, ".", ",", "NT$", "$1", 12345, "12,345"},
		{0, ".", ",", "NT$", "$1", 123456, "123,456"},
		{0, ".", ",", "NT$", "$1", 1234567, "1,234,567"},
		{0, ".", ",", "NT$", "$1", 12345678, "12,345,678"},
		{0, ".", ",", "NT$", "$1", 123456789, "123,456,789"},

		{0, ".", ",", "NT$", "$1", -1, "-1"},
		{0, ".", ",", "NT$", "$1", -12, "-12"},
		{0, ".", ",", "NT$", "$1", -123, "-123"},
		{0, ".", ",", "NT$", "$1", -1234, "-1,234"},
		{0, ".", ",", "NT$", "$1", -12345, "-12,345"},
		{0, ".", ",", "NT$", "$1", -123456, "-123,456"},
		{0, ".", ",", "NT$", "$1", -1234567, "-1,234,567"},
		{0, ".", ",", "NT$", "$1", -12345678, "-12,345,678"},
		{0, ".", ",", "NT$", "$1", -123456789, "-123,456,789"},
	}

	for _, tc := range tcs {
		formatter := NewFormatter(tc.fraction, tc.decimal, tc.thousand, tc.grapheme, tc.template)
		r := formatter.FormatAmount(tc.amount)

		if r != tc.expected {
			t.Errorf("Expected %d formatted to be %s got %s", tc.amount, tc.expected, r)
		}
	}
}

func TestFormatter_ToMajorUnits(t *testing.T) {
	tcs := []struct {
		fraction int
		decimal  string
		thousand string
		grapheme string
		template string
		amount   int64
		expected float64
	}{
		{2, ".", ",", "$", "1 $", 0, 0.00},
		{2, ".", ",", "$", "1 $", 1, 0.01},
		{2, ".", ",", "$", "1 $", 12, 0.12},
		{2, ".", ",", "$", "1 $", 123, 1.23},
		{2, ".", ",", "$", "1 $", 1234, 12.34},
		{2, ".", ",", "$", "1 $", 12345, 123.45},
		{2, ".", ",", "$", "1 $", 123456, 1234.56},
		{2, ".", ",", "$", "1 $", 1234567, 12345.67},
		{2, ".", ",", "$", "1 $", 12345678, 123456.78},
		{2, ".", ",", "$", "1 $", 123456789, 1234567.89},

		{2, ".", ",", "$", "1 $", -1, -0.01},
		{2, ".", ",", "$", "1 $", -12, -0.12},
		{2, ".", ",", "$", "1 $", -123, -1.23},
		{2, ".", ",", "$", "1 $", -1234, -12.34},
		{2, ".", ",", "$", "1 $", -12345, -123.45},
		{2, ".", ",", "$", "1 $", -123456, -1234.56},
		{2, ".", ",", "$", "1 $", -1234567, -12345.67},
		{2, ".", ",", "$", "1 $", -12345678, -123456.78},
		{2, ".", ",", "$", "1 $", -123456789, -1234567.89},

		{3, ".", "", "$", "1 $", 1, 0.001},
		{3, ".", "", "$", "1 $", 12, 0.012},
		{3, ".", "", "$", "1 $", 123, 0.123},
		{3, ".", "", "$", "1 $", 1234, 1.234},
		{3, ".", "", "$", "1 $", 12345, 12.345},
		{3, ".", "", "$", "1 $", 123456, 123.456},
		{3, ".", "", "$", "1 $", 1234567, 1234.567},
		{3, ".", "", "$", "1 $", 12345678, 12345.678},
		{3, ".", "", "$", "1 $", 123456789, 123456.789},

		{2, ".", ",", "£", "$1", 1, 0.01},
		{2, ".", ",", "£", "$1", 12, 0.12},
		{2, ".", ",", "£", "$1", 123, 1.23},
		{2, ".", ",", "£", "$1", 1234, 12.34},
		{2, ".", ",", "£", "$1", 12345, 123.45},
		{2, ".", ",", "£", "$1", 123456, 1234.56},
		{2, ".", ",", "£", "$1", 1234567, 12345.67},
		{2, ".", ",", "£", "$1", 12345678, 123456.78},
		{2, ".", ",", "£", "$1", 123456789, 1234567.89},

		{0, ".", ",", "NT$", "$1", 1, 1},
		{0, ".", ",", "NT$", "$1", 12, 12},
		{0, ".", ",", "NT$", "$1", 123, 123},
		{0, ".", ",", "NT$", "$1", 1234, 1234},
		{0, ".", ",", "NT$", "$1", 12345, 12345},
		{0, ".", ",", "NT$", "$1", 123456, 123456},
		{0, ".", ",", "NT$", "$1", 1234567, 1234567},
		{0, ".", ",", "NT$", "$1", 12345678, 12345678},
		{0, ".", ",", "NT$", "$1", 123456789, 123456789},

		{0, ".", ",", "NT$", "$1", -1, -1},
		{0, ".", ",", "NT$", "$1", -12, -12},
		{0, ".", ",", "NT$", "$1", -123, -123},
		{0, ".", ",", "NT$", "$1", -1234, -1234},
		{0, ".", ",", "NT$", "$1", -12345, -12345},
		{0, ".", ",", "NT$", "$1", -123456, -123456},
		{0, ".", ",", "NT$", "$1", -1234567, -1234567},
		{0, ".", ",", "NT$", "$1", -12345678, -12345678},
		{0, ".", ",", "NT$", "$1", -123456789, -123456789},
	}

	for _, tc := range tcs {
		formatter := NewFormatter(tc.fraction, tc.decimal, tc.thousand, tc.grapheme, tc.template)
		r := formatter.ToMajorUnits(tc.amount)

		if r != tc.expected {
			t.Errorf("Expected %d formatted to major units to be %f got %f", tc.amount, tc.expected, r)
		}
	}
}
