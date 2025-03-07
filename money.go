package money

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Injection points for backward compatibility.
// If you need to keep your JSON marshal/unmarshal way, overwrite them like below.
//
//	money.UnmarshalJSON = func (m *Money, b []byte) error { ... }
//	money.MarshalJSON = func (m Money) ([]byte, error) { ... }
var (
	// UnmarshalJSON is injection point of json.Unmarshaller for money.Money
	UnmarshalJSON = unmarshalJSON
	// MarshalJSON is injection point of json.Marshaller for money.Money
	MarshalJSON = marshalJSON

	// ErrCurrencyMismatch happens when two compared Money don't have the same currency.
	ErrCurrencyMismatch = errors.New("currencies don't match")

	// ErrInvalidJSON happens when the default money.UnmarshalJSON fails to unmarshal Money because of invalid data.
	ErrInvalidJSON = errors.New("invalid json")
)

func unmarshalJSON(m *Money, b []byte) error {
	data := make(map[string]interface{})
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	var amount string
	if amountRaw, ok := data["amount"]; ok {
		amount, ok = amountRaw.(string)
		if !ok {
			return ErrInvalidJSON
		}
	}

	var currency string
	if currencyRaw, ok := data["currency"]; ok {
		currency, ok = currencyRaw.(string)
		if !ok {
			return ErrInvalidJSON
		}
	}

	var ref *Money
	if amount == "" && currency == "" {
		ref = &Money{}
	} else {
		ref, err = NewFromString(amount, currency)
	}

	if err != nil {
		return err
	}

	*m = *ref
	return nil
}

func marshalJSON(m Money) ([]byte, error) {
	if m == (Money{}) {
		m = Money{0, newCurrency("").get()}
	}

	buff := bytes.NewBufferString(fmt.Sprintf(`{"amount": "%s", "currency": "%s"}`, m.Amount(), m.CurrencyCode()))
	return buff.Bytes(), nil
}

// Amount is a data structure that stores the amount being used for calculations.
type Amount = int64

// Money represents monetary value information, stores
// currency and amount value.
type Money struct {
	amount   Amount
	currency *Currency
}

// New creates and returns new instance of Money.
func New(amount int64, currencyCode string) (*Money, error) {
	currency := GetCurrency(currencyCode)
	if currency == nil {
		return nil, fmt.Errorf("invalid currency '%s'", currencyCode)
	}

	return &Money{
		amount:   amount,
		currency: currency,
	}, nil
}

// NewFromFloat creates and returns new instance of Money from a float64.
// Always rounding trailing decimals down.
//
// Using NewFromFloat WILL CAUSE PRECISION ISSUES IN CERTAIN CASES!
//
// For example:
//
//	m := money.NewFromFloat(1.15, money.USD)
//	fmt.Println(m.Amount())
//
// The above code will output 114 instead of 115.
func NewFromFloat(amount float64, currencyCode string) (*Money, error) {
	currency := GetCurrency(currencyCode)
	if currency == nil {
		return nil, fmt.Errorf("invalid currency '%s'", currencyCode)
	}
	currencyDecimals := math.Pow10(currency.Fraction)

	return &Money{
		amount:   int64(amount * currencyDecimals),
		currency: currency,
	}, nil
}

// NewFromString creates and returns new instance of Money from a string.
// Can only parse simple float-like strings, like "1.23" USD or "1.5" ARS, not "1.23 USD", "$1.23" or "1.000" USD.
func NewFromString(amount string, currencyCode string) (*Money, error) {
	currency := GetCurrency(currencyCode)
	if currency == nil {
		return nil, fmt.Errorf("invalid currency '%s'", currencyCode)
	}

	fraction := currency.Fraction

	toParse := amount
	var decimals int
	if pointIndex := strings.Index(amount, currency.Decimal); pointIndex != -1 {
		decimals = len(amount) - pointIndex - 1
		if decimals > fraction {
			decimals = fraction
		}
		toParse = amount[:pointIndex] + amount[pointIndex+1:pointIndex+1+decimals]
	}

	parsed, err := strconv.ParseInt(toParse, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amount '%s'", amount)
	}

	for d := decimals; d < fraction; d++ {
		parsed *= 10
	}

	return &Money{
		amount:   parsed,
		currency: currency,
	}, nil
}

// Currency returns the currency used by Money.
func (m *Money) CurrencyCode() string {
	return m.currency.Code
}

// AmountUnformatted returns a copy of the internal monetary value as an int64.
func (m *Money) AmountUnformatted() int64 {
	return m.amount
}

func (m *Money) Amount() string {
	currency := m.currency.get()
	return currency.Formatter().FormatAmount(m.amount)
}

// SameCurrency check if given Money is equals by currency.
func (m *Money) SameCurrency(om *Money) bool {
	return m.currency.equals(om.currency)
}

func (m *Money) assertSameCurrency(om *Money) error {
	if !m.SameCurrency(om) {
		return ErrCurrencyMismatch
	}

	return nil
}

func (m *Money) compare(om *Money) int {
	switch {
	case m.amount > om.amount:
		return 1
	case m.amount < om.amount:
		return -1
	}

	return 0
}

// Equals checks equality between two Money types.
func (m *Money) Equals(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 0, nil
}

// GreaterThan checks whether the value of Money is greater than the other.
func (m *Money) GreaterThan(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 1, nil
}

// GreaterThanOrEqual checks whether the value of Money is greater or equal than the other.
func (m *Money) GreaterThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) >= 0, nil
}

// LessThan checks whether the value of Money is less than the other.
func (m *Money) LessThan(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == -1, nil
}

// LessThanOrEqual checks whether the value of Money is less or equal than the other.
func (m *Money) LessThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) <= 0, nil
}

// IsZero returns boolean of whether the value of Money is equals to zero.
func (m *Money) IsZero() bool {
	return m.amount == 0
}

// IsPositive returns boolean of whether the value of Money is positive.
func (m *Money) IsPositive() bool {
	return m.amount > 0
}

// IsNegative returns boolean of whether the value of Money is negative.
func (m *Money) IsNegative() bool {
	return m.amount < 0
}

// Absolute returns new Money struct from given Money using absolute monetary value.
func (m *Money) Absolute() *Money {
	return &Money{amount: mutate.calc.absolute(m.amount), currency: m.currency}
}

// Negative returns new Money struct from given Money using negative monetary value.
func (m *Money) Negative() *Money {
	return &Money{amount: mutate.calc.negative(m.amount), currency: m.currency}
}

// Add returns new Money struct with value representing sum of Self and Other Money.
func (m *Money) Add(om *Money) (*Money, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return nil, err
	}

	return &Money{amount: mutate.calc.add(m.amount, om.amount), currency: m.currency}, nil
}

// Subtract returns new Money struct with value representing difference of Self and Other Money.
func (m *Money) Subtract(om *Money) (*Money, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return nil, err
	}

	return &Money{amount: mutate.calc.subtract(m.amount, om.amount), currency: m.currency}, nil
}

// Multiply returns new Money struct with value representing Self multiplied value by multiplier.
func (m *Money) Multiply(mul int64) *Money {
	return &Money{amount: mutate.calc.multiply(m.amount, mul), currency: m.currency}
}

// Round returns new Money struct with value rounded to nearest zero.
func (m *Money) Round() *Money {
	return &Money{amount: mutate.calc.round(m.amount, m.currency.Fraction), currency: m.currency}
}

// Split returns slice of Money structs with split Self value in given number.
// After division leftover pennies will be distributed round-robin amongst the parties.
// This means that parties listed first will likely receive more pennies than ones that are listed later.
func (m *Money) Split(n int) ([]*Money, error) {
	if n <= 0 {
		return nil, errors.New("split must be higher than zero")
	}

	a := mutate.calc.divide(m.amount, int64(n))
	ms := make([]*Money, n)

	for i := 0; i < n; i++ {
		ms[i] = &Money{amount: a, currency: m.currency}
	}

	r := mutate.calc.modulus(m.amount, int64(n))
	l := mutate.calc.absolute(r)
	// Add leftovers to the first parties.

	v := int64(1)
	if m.amount < 0 {
		v = -1
	}
	for p := 0; l != 0; p++ {
		ms[p].amount = mutate.calc.add(ms[p].amount, v)
		l--
	}

	return ms, nil
}

// Allocate returns slice of Money structs with split Self value in given ratios.
// It lets split money by given ratios without losing pennies and as Split operations distributes
// leftover pennies amongst the parties with round-robin principle.
func (m *Money) Allocate(rs ...int) ([]*Money, error) {
	if len(rs) == 0 {
		return nil, errors.New("no ratios specified")
	}

	// Calculate sum of ratios.
	var sum uint
	for _, r := range rs {
		if r < 0 {
			return nil, errors.New("negative ratios not allowed")
		}
		sum += uint(r)
	}

	var total int64
	ms := make([]*Money, 0, len(rs))
	for _, r := range rs {
		party := &Money{
			amount:   mutate.calc.allocate(m.amount, uint(r), sum),
			currency: m.currency,
		}

		ms = append(ms, party)
		total += party.amount
	}

	// if the sum of all ratios is zero, then we just returns zeros and don't do anything
	// with the leftover
	if sum == 0 {
		return ms, nil
	}

	// Calculate leftover value and divide to first parties.
	lo := m.amount - total
	sub := int64(1)
	if lo < 0 {
		sub = -sub
	}

	for p := 0; lo != 0; p++ {
		ms[p].amount = mutate.calc.add(ms[p].amount, sub)
		lo -= sub
	}

	return ms, nil
}

// Display lets represent Money struct as string in given Currency value.
func (m *Money) Display() string {
	c := m.currency.get()
	return c.Formatter().Format(m.amount)
}

// AsMajorUnits lets represent Money struct as subunits (float64) in given Currency value
func (m *Money) AsMajorUnits() float64 {
	c := m.currency.get()
	return c.Formatter().ToMajorUnits(m.amount)
}

// UnmarshalJSON is implementation of json.Unmarshaller
func (m *Money) UnmarshalJSON(b []byte) error {
	return UnmarshalJSON(m, b)
}

// MarshalJSON is implementation of json.Marshaller
func (m Money) MarshalJSON() ([]byte, error) {
	return MarshalJSON(m)
}

// Compare function compares two money of the same type
//
//	if m.amount > om.amount returns (1, nil)
//	if m.amount == om.amount returns (0, nil
//	if m.amount < om.amount returns (-1, nil)
//
// If compare moneys from distinct currency, return (m.amount, ErrCurrencyMismatch)
func (m *Money) Compare(om *Money) (int, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return int(m.amount), err
	}

	return m.compare(om), nil
}
