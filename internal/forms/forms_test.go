package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("Got invalid but should have got valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Form shows valid whem required field missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	has := form.Has("whatever")
	if has == true {
		t.Error("form shows has field when it does not exist")
	}

	postedData := url.Values{}

	postedData.Add("a", "a")
	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("shows form does not has a field when it does")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("Form shows minlength for a non-existent field")
	}

	isError := form.Errors.Get("x")

	if isError == "" {
		t.Error("should have an error but did not get one")
	}

	postedValues = url.Values{}
	postedValues.Add("some_field", "some value")

	form = New(postedValues)
	form.MinLength("some_field", 100)

	if form.Valid() {
		t.Error("shows minlength of 100 met when passed data is shorter")
	}

	postedValues = url.Values{}
	postedValues.Add("another_field", "darshan")

	form = New(postedValues)
	form.MinLength("another_field", 2)

	if !form.Valid() {
		t.Error("shows minlength not met when it should be met")
	}

	isError = form.Errors.Get("another_field")
	if isError != "" {
		t.Error("Should not have an error but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	form.IsEmail("x")

	if form.Valid() {
		t.Error("shows valid email for a non-existent field")
	}

	postedValues.Add("email", "jejf")

	form = New(postedValues)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("shows valid email for a non-valid email")
	}

	postedValues.Add("another_email", "a@a.com")
	form = New(postedValues)

	form.IsEmail("another_email")
	if !form.Valid() {
		t.Error("shows invalid email for a valid email")
	}
}
