package money

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	m, _ := New(1, EUR)

	if m.amount != 1 {
		t.Errorf("Expected %d got %d", 1, m.amount)
	}

	if m.currency.Code != EUR {
		t.Errorf("Expected currency %s got %s", EUR, m.currency.Code)
	}

	m, _ = New(-100, EUR)

	if m.amount != -100 {
		t.Errorf("Expected %d got %d", -100, m.amount)
	}
}

func TestCurrency(t *testing.T) {
	code := "MOCK"
	decimals := 5
	AddCurrency(code, "M$", "1 $", ".", ",", decimals)
	m, _ := New(1, code)

	currency := GetCurrency(m.CurrencyCode())

	c := currency.Code
	if c != code {
		t.Errorf("Expected %s got %s", code, c)
	}
	f := currency.Fraction
	if f != decimals {
		t.Errorf("Expected %d got %d", decimals, f)
	}
}

func TestMoney_SameCurrency(t *testing.T) {
	m, _ := New(0, EUR)
	om, _ := New(0, USD)

	if m.SameCurrency(om) {
		t.Errorf("Expected %s not to be same as %s", m.currency.Code, om.currency.Code)
	}

	om, _ = New(0, EUR)

	if !m.SameCurrency(om) {
		t.Errorf("Expected %s to be same as %s", m.currency.Code, om.currency.Code)
	}
}

func TestMoney_Equals(t *testing.T) {
	m, _ := New(0, EUR)
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, false},
		{0, true},
		{1, false},
	}

	for _, tc := range tcs {
		om, _ := New(tc.amount, EUR)
		r, err := m.Equals(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Equals %d == %t got %t", m.amount,
				om.amount, tc.expected, r)
		}
	}
}

func TestMoney_Equals_DifferentCurrencies(t *testing.T) {
	t.Parallel()

	eur, _ := New(0, EUR)
	usd, _ := New(0, USD)

	_, err := eur.Equals(usd)
	if err == nil || !errors.Is(ErrCurrencyMismatch, err) {
		t.Errorf("Expected Equals to return %q, got %v", ErrCurrencyMismatch.Error(), err)
	}
}

func TestMoney_GreaterThan(t *testing.T) {
	m, _ := New(0, EUR)
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, true},
		{0, false},
		{1, false},
	}

	for _, tc := range tcs {
		om, _ := New(tc.amount, EUR)
		r, err := m.GreaterThan(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Greater Than %d == %t got %t", m.amount,
				om.amount, tc.expected, r)
		}
	}
}

func TestMoney_GreaterThanOrEqual(t *testing.T) {
	m, _ := New(0, EUR)
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, true},
		{0, true},
		{1, false},
	}

	for _, tc := range tcs {
		om, _ := New(tc.amount, EUR)
		r, err := m.GreaterThanOrEqual(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Equals Or Greater Than %d == %t got %t", m.amount,
				om.amount, tc.expected, r)
		}
	}
}

func TestMoney_LessThan(t *testing.T) {
	m, _ := New(0, EUR)
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, false},
		{0, false},
		{1, true},
	}

	for _, tc := range tcs {
		om, _ := New(tc.amount, EUR)
		r, err := m.LessThan(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Less Than %d == %t got %t", m.amount,
				om.amount, tc.expected, r)
		}
	}
}

func TestMoney_LessThanOrEqual(t *testing.T) {
	m, _ := New(0, EUR)
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, false},
		{0, true},
		{1, true},
	}

	for _, tc := range tcs {
		om, _ := New(tc.amount, EUR)
		r, err := m.LessThanOrEqual(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Equal Or Less Than %d == %t got %t", m.amount,
				om.amount, tc.expected, r)
		}
	}
}

func TestMoney_IsZero(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, false},
		{0, true},
		{1, false},
	}

	for _, tc := range tcs {
		m, _ := New(tc.amount, EUR)
		r := m.IsZero()

		if r != tc.expected {
			t.Errorf("Expected %d to be zero == %t got %t", m.amount, tc.expected, r)
		}
	}
}

func TestMoney_IsNegative(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, true},
		{0, false},
		{1, false},
	}

	for _, tc := range tcs {
		m, _ := New(tc.amount, EUR)
		r := m.IsNegative()

		if r != tc.expected {
			t.Errorf("Expected %d to be negative == %t got %t", m.amount,
				tc.expected, r)
		}
	}
}

func TestMoney_IsPositive(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, false},
		{0, false},
		{1, true},
	}

	for _, tc := range tcs {
		m, _ := New(tc.amount, EUR)
		r := m.IsPositive()

		if r != tc.expected {
			t.Errorf("Expected %d to be positive == %t got %t", m.amount,
				tc.expected, r)
		}
	}
}

func TestMoney_Absolute(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected int64
	}{
		{-1, 1},
		{0, 0},
		{1, 1},
	}

	for _, tc := range tcs {
		m, _ := New(tc.amount, EUR)
		r := m.Absolute().amount

		if r != tc.expected {
			t.Errorf("Expected absolute %d to be %d got %d", m.amount,
				tc.expected, r)
		}
	}
}

func TestMoney_Negative(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected int64
	}{
		{-1, -1},
		{0, -0},
		{1, -1},
	}

	for _, tc := range tcs {
		m, _ := New(tc.amount, EUR)
		r := m.Negative().amount

		if r != tc.expected {
			t.Errorf("Expected absolute %d to be %d got %d", m.amount,
				tc.expected, r)
		}
	}
}

func TestMoney_Add(t *testing.T) {
	tcs := []struct {
		amount1  int64
		amount2  int64
		expected int64
	}{
		{5, 5, 10},
		{10, 5, 15},
		{1, -1, 0},
	}

	for _, tc := range tcs {
		m, _ := New(tc.amount1, EUR)
		om, _ := New(tc.amount2, EUR)
		r, err := m.Add(om)

		if err != nil {
			t.Error(err)
		}

		if r.AmountUnformatted() != tc.expected {
			t.Errorf("Expected %d + %d = %d got %d", tc.amount1, tc.amount2,
				tc.expected, r.amount)
		}
	}
}

func TestMoney_Add2(t *testing.T) {
	m, _ := New(100, EUR)
	dm, _ := New(100, GBP)
	r, err := m.Add(dm)

	if r != nil || err == nil {
		t.Error("Expected err")
	}
}

func TestMoney_Subtract(t *testing.T) {
	tcs := []struct {
		amount1  int64
		amount2  int64
		expected int64
	}{
		{5, 5, 0},
		{10, 5, 5},
		{1, -1, 2},
	}

	for _, tc := range tcs {
		m, _ := New(tc.amount1, EUR)
		om, _ := New(tc.amount2, EUR)
		r, err := m.Subtract(om)

		if err != nil {
			t.Error(err)
		}

		if r.amount != tc.expected {
			t.Errorf("Expected %d - %d = %d got %d", tc.amount1, tc.amount2,
				tc.expected, r.amount)
		}
	}
}

func TestMoney_Subtract2(t *testing.T) {
	m, _ := New(100, EUR)
	dm, _ := New(100, GBP)
	r, err := m.Subtract(dm)

	if r != nil || err == nil {
		t.Error("Expected err")
	}
}

func TestMoney_Multiply(t *testing.T) {
	tcs := []struct {
		amount     int64
		multiplier int64
		expected   int64
	}{
		{5, 5, 25},
		{10, 5, 50},
		{1, -1, -1},
		{1, 0, 0},
	}

	for _, tc := range tcs {
		m, _ := New(tc.amount, EUR)
		r := m.Multiply(tc.multiplier).amount

		if r != tc.expected {
			t.Errorf("Expected %d * %d = %d got %d", tc.amount, tc.multiplier, tc.expected, r)
		}
	}
}

func TestMoney_Round(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected int64
	}{
		{125, 100},
		{175, 200},
		{349, 300},
		{351, 400},
		{0, 0},
		{-1, 0},
		{-75, -100},
	}

	for _, tc := range tcs {
		m, _ := New(tc.amount, EUR)
		r := m.Round().amount

		if r != tc.expected {
			t.Errorf("Expected rounded %d to be %d got %d", tc.amount, tc.expected, r)
		}
	}
}

func TestMoney_RoundWithExponential(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected int64
	}{
		{12555, 13000},
	}

	for _, tc := range tcs {
		AddCurrency("CUR", "*", "$1", ".", ",", 3)
		m, _ := New(tc.amount, "CUR")
		r := m.Round().amount

		if r != tc.expected {
			t.Errorf("Expected rounded %d to be %d got %d", tc.amount, tc.expected, r)
		}
	}
}

func TestMoney_Split(t *testing.T) {
	tcs := []struct {
		amount   int64
		split    int
		expected []int64
	}{
		{100, 3, []int64{34, 33, 33}},
		{100, 4, []int64{25, 25, 25, 25}},
		{5, 3, []int64{2, 2, 1}},
		{-101, 4, []int64{-26, -25, -25, -25}},
		{-101, 4, []int64{-26, -25, -25, -25}},
		{-2, 3, []int64{-1, -1, 0}},
	}

	for _, tc := range tcs {
		m, _ := New(tc.amount, EUR)
		var rs []int64
		split, _ := m.Split(tc.split)

		for _, party := range split {
			rs = append(rs, party.amount)
		}

		if !reflect.DeepEqual(tc.expected, rs) {
			t.Errorf("Expected split of %d to be %v got %v", tc.amount, tc.expected, rs)
		}
	}
}

func TestMoney_Split2(t *testing.T) {
	m, _ := New(100, EUR)
	r, err := m.Split(-10)

	if r != nil || err == nil {
		t.Error("Expected err")
	}
}

func TestMoney_Allocate(t *testing.T) {
	tcs := []struct {
		amount   int64
		ratios   []int
		expected []int64
	}{
		{100, []int{50, 50}, []int64{50, 50}},
		{100, []int{30, 30, 30}, []int64{34, 33, 33}},
		{200, []int{25, 25, 50}, []int64{50, 50, 100}},
		{5, []int{50, 25, 25}, []int64{3, 1, 1}},
		{0, []int{0, 0, 0, 0}, []int64{0, 0, 0, 0}},
		{0, []int{50, 10}, []int64{0, 0}},
		{10, []int{0, 100}, []int64{0, 10}},
		{10, []int{0, 0}, []int64{0, 0}},
	}

	for _, tc := range tcs {
		m, _ := New(tc.amount, EUR)
		var rs []int64
		split, _ := m.Allocate(tc.ratios...)

		for _, party := range split {
			rs = append(rs, party.amount)
		}

		if !reflect.DeepEqual(tc.expected, rs) {
			t.Errorf("Expected allocation of %d for ratios %v to be %v got %v", tc.amount, tc.ratios,
				tc.expected, rs)
		}
	}
}

func TestMoney_Allocate2(t *testing.T) {
	m, _ := New(100, EUR)
	r, err := m.Allocate()

	if r != nil || err == nil {
		t.Error("Expected err")
	}
}

func TestMoney_Format(t *testing.T) {
	tcs := []struct {
		amount   int64
		code     string
		expected string
	}{
		{100, GBP, "£1.00"},
	}

	for _, tc := range tcs {
		m, _ := New(tc.amount, tc.code)
		r := m.Display()

		if r != tc.expected {
			t.Errorf("Expected formatted %d to be %s got %s", tc.amount, tc.expected, r)
		}
	}
}

func TestMoney_Display(t *testing.T) {
	tcs := []struct {
		amount   int64
		code     string
		expected string
	}{
		{100, AED, "1.00 .\u062f.\u0625"},
		{1, USD, "$0.01"},
	}

	for _, tc := range tcs {
		m, _ := New(tc.amount, tc.code)
		r := m.Display()

		if r != tc.expected {
			t.Errorf("Expected formatted %d to be %s got %s", tc.amount, tc.expected, r)
		}
	}
}

func TestMoney_AsMajorUnits(t *testing.T) {
	tcs := []struct {
		amount   int64
		code     string
		expected float64
	}{
		{100, AED, 1.00},
		{1, USD, 0.01},
	}

	for _, tc := range tcs {
		m, _ := New(tc.amount, tc.code)
		r := m.AsMajorUnits()

		if r != tc.expected {
			t.Errorf("Expected value as major units of %d to be %f got %f", tc.amount, tc.expected, r)
		}
	}
}

func TestMoney_Allocate3(t *testing.T) {
	pound, _ := New(100, GBP)
	parties, err := pound.Allocate(33, 33, 33)

	if err != nil {
		t.Error(err)
	}

	if parties[0].Display() != "£0.34" {
		t.Errorf("Expected %s got %s", "£0.34", parties[0].Display())
	}

	if parties[1].Display() != "£0.33" {
		t.Errorf("Expected %s got %s", "£0.33", parties[1].Display())
	}

	if parties[2].Display() != "£0.33" {
		t.Errorf("Expected %s got %s", "£0.33", parties[2].Display())
	}
}

func TestMoney_Comparison(t *testing.T) {
	pound, _ := New(100, GBP)
	twoPounds, _ := New(200, GBP)
	twoEuros, _ := New(200, EUR)

	if r, err := pound.GreaterThan(twoPounds); err != nil || r {
		t.Errorf("Expected %d Greater Than %d == %t got %t", pound.amount,
			twoPounds.amount, false, r)
	}

	if r, err := pound.LessThan(twoPounds); err != nil || !r {
		t.Errorf("Expected %d Less Than %d == %t got %t", pound.amount,
			twoPounds.amount, true, r)
	}

	if r, err := pound.LessThan(twoEuros); err == nil || r {
		t.Error("Expected err")
	}

	if r, err := pound.GreaterThan(twoEuros); err == nil || r {
		t.Error("Expected err")
	}

	if r, err := pound.Equals(twoEuros); err == nil || r {
		t.Error("Expected err")
	}

	if r, err := pound.LessThanOrEqual(twoEuros); err == nil || r {
		t.Error("Expected err")
	}

	if r, err := pound.GreaterThanOrEqual(twoEuros); err == nil || r {
		t.Error("Expected err")
	}

	if r, err := twoPounds.Compare(pound); r != 1 && err != nil {
		t.Errorf("Expected %d Greater Than %d == %d got %d", pound.amount,
			twoPounds.amount, 1, r)
	}

	if r, err := pound.Compare(twoPounds); r != -1 && err != nil {
		t.Errorf("Expected %d Less Than %d == %d got %d", pound.amount,
			twoPounds.amount, -1, r)
	}

	if _, err := pound.Compare(twoEuros); err != ErrCurrencyMismatch {
		t.Error("Expected err")
	}

	anotherTwoEuros, _ := New(200, EUR)
	if r, err := twoEuros.Compare(anotherTwoEuros); r != 0 && err != nil {
		t.Errorf("Expected %d Equals to %d == %d got %d", anotherTwoEuros.amount,
			twoEuros.amount, 0, r)
	}

}

func TestMoney_Currency(t *testing.T) {
	pound, _ := New(100, GBP)

	if pound.CurrencyCode() != GBP {
		t.Errorf("Expected %s got %s", GBP, pound.CurrencyCode())
	}
}

func TestMoney_Amount(t *testing.T) {
	pound, _ := New(100, GBP)

	if pound.AmountUnformatted() != 100 {
		t.Errorf("Expected %d got %d", 100, pound.AmountUnformatted())
	}
}

func TestNewFromFloat(t *testing.T) {
	m, _ := NewFromFloat(12.34, EUR)

	if m.amount != 1234 {
		t.Errorf("Expected %d got %d", 1234, m.amount)
	}

	if m.currency.Code != EUR {
		t.Errorf("Expected currency %s got %s", EUR, m.currency.Code)
	}

	m, _ = NewFromFloat(-0.125, EUR)

	if m.amount != -12 {
		t.Errorf("Expected %d got %d", -12, m.amount)
	}
}

func TestNewFromString(t *testing.T) {
	m, err := NewFromString("12.34", EUR)

	if err != nil {
		t.Error(err)
	}

	if m.amount != 1234 {
		t.Errorf("Expected %d got %d", 1234, m.amount)
	}

	if m.currency.Code != EUR {
		t.Errorf("Expected currency %s got %s", EUR, m.currency.Code)
	}

	m, err = NewFromString("-1.12345", EUR)

	if err != nil {
		t.Error(err)
	}

	if m.amount != -112 {
		t.Errorf("Expected %d got %d", -112, m.amount)
	}

	_, err = NewFromString("invalid_input", EUR)

	if err.Error() != "invalid amount 'invalid_input'" {
		t.Error(err)
	}
}

func TestDefaultMarshal(t *testing.T) {
	given, _ := New(12345, IQD)
	expected := `{"amount":"12.345","currency":"IQD"}`

	b, err := json.Marshal(given)

	if err != nil {
		t.Error(err)
	}

	if string(b) != expected {
		t.Errorf("Expected %s got %s", expected, string(b))
	}

	given = &Money{}
	expected = `{"amount":"0.00","currency":""}`

	b, err = json.Marshal(given)

	if err != nil {
		t.Error(err)
	}

	if string(b) != expected {
		t.Errorf("Expected %s got %s", expected, string(b))
	}
}

func TestCustomMarshal(t *testing.T) {
	given, _ := New(12345, IQD)
	expected := `{"amount":12345,"currency_code":"IQD","currency_fraction":3}`

	MarshalJSON = func(m Money) ([]byte, error) {
		currency := GetCurrency(m.CurrencyCode())

		buff := bytes.NewBufferString(fmt.Sprintf(`{"amount": %d, "currency_code": "%s", "currency_fraction": %d}`, m.AmountUnformatted(), currency.Code, currency.Fraction))
		return buff.Bytes(), nil
	}

	b, err := json.Marshal(given)

	if err != nil {
		t.Error(err)
	}

	if string(b) != expected {
		t.Errorf("Expected %s got %s", expected, string(b))
	}
}

func TestDefaultUnmarshal(t *testing.T) {
	given := `{"amount": "100.12", "currency":"USD"}`
	expected := "$100.12"
	var m Money
	err := json.Unmarshal([]byte(given), &m)
	if err != nil {
		t.Error(err)
	}

	if m.Display() != expected {
		t.Errorf("Expected %s got %s", expected, m.Display())
	}

	given = `{}`
	err = json.Unmarshal([]byte(given), &m)
	if err != nil {
		t.Error(err)
	}

	if m != (Money{}) {
		t.Errorf("Expected zero value, got %+v", m)
	}

	given = `{"amount": "foo", "currency": "USD"}`
	err = json.Unmarshal([]byte(given), &m)
	expectedErr := fmt.Errorf("invalid amount '%s'", "foo")
	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected `%v`, got `%+v`", expectedErr, err)
	}

	given = `{"amount": 1234, "currency": 1234}`
	err = json.Unmarshal([]byte(given), &m)
	if !errors.Is(err, ErrInvalidJSON) {
		t.Errorf("Expected ErrInvalidJSON, got %+v", err)
	}
}

func TestCustomUnmarshal(t *testing.T) {
	given := `{"amount": 10012, "currency_code":"USD", "currency_fraction":2}`
	expected := "$100.12"
	UnmarshalJSON = func(m *Money, b []byte) error {
		data := make(map[string]interface{})
		err := json.Unmarshal(b, &data)
		if err != nil {
			return err
		}
		ref, _ := New(int64(data["amount"].(float64)), data["currency_code"].(string))
		*m = *ref
		return nil
	}

	var m Money
	err := json.Unmarshal([]byte(given), &m)
	if err != nil {
		t.Error(err)
	}

	if m.Display() != expected {
		t.Errorf("Expected %s got %s", expected, m.Display())
	}
}
