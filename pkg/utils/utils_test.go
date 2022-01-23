package utils

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"testing"
)

func TestHumanizeString(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"API", "API"},
		{"OrderID", "Order ID"},
		{"OrderItem", "Order Item"},
		{"orderItem", "Order Item"},
		{"OrderIDItem", "Order ID Item"},
		{"OrderItemID", "Order Item ID"},
		{"VIEW SITE", "VIEW SITE"},
		{"Order Item", "Order Item"},
		{"Order ITEM", "Order ITEM"},
		{"ORDER Item", "ORDER Item"},
	}
	for _, c := range cases {
		if got := HumanizeString(c.input); got != c.want {
			t.Errorf("HumanizeString(%q) = %q; want %q", c.input, got, c.want)
		}
	}
}

func TestToParamString(t *testing.T) {
	results := map[string]string{
		"OrderItem":  "order_item",
		"order item": "order_item",
		"Order Item": "order_item",
		"FAQ":        "faq",
		"FAQPage":    "faq_page",
		"!help_id":   "!help_id",
		"help id":    "help_id",
		"语言":         "yu-yan",
	}

	for key, value := range results {
		if ToParamString(key) != value {
			t.Errorf("%v to params should be %v, but got %v", key, value, ToParamString(key))
		}
	}
}

func TestPatchURL(t *testing.T) {
	var cases = []struct {
		original string
		input    []interface{}
		want     string
		err      error
	}{
		{
			original: "http://app.bhojpur.net/admin/orders?locale=global&q=dotnet&test=1#test",
			input:    []interface{}{"locale", "cn"},
			want:     "http://app.bhojpur.net/admin/orders?locale=cn&q=dotnet&test=1#test",
		},
		{
			original: "http://app.bhojpur.net/admin/orders?locale=global&q=dotnet&test=1#test",
			input:    []interface{}{"locale", ""},
			want:     "http://app.bhojpur.net/admin/orders?q=dotnet&test=1#test",
		},
	}
	for _, c := range cases {
		// u, _ := url.Parse(c.original)
		// context := Context{Context: &bhojpur.Context{Request: &http.Request{URL: u}}}
		got, err := PatchURL(c.original, c.input...)
		if c.err != nil {
			if err == nil || err.Error() != c.err.Error() {
				t.Errorf("got error %s; want %s", err, c.err)
			}
		} else {
			if err != nil {
				t.Error(err)
			}
			if got != c.want {
				t.Errorf("context.PatchURL = %s; c.want %s", got, c.want)
			}
		}
	}
}

func TestJoinURL(t *testing.T) {
	var cases = []struct {
		original string
		input    []interface{}
		want     string
		err      error
	}{
		{
			original: "http://app.bhojpur.net",
			input:    []interface{}{"admin"},
			want:     "http://app.bhojpur.net/admin",
		},
		{
			original: "http://app.bhojpur.net",
			input:    []interface{}{"/admin"},
			want:     "http://app.bhojpur.net/admin",
		},
		{
			original: "http://app.bhojpur.net/",
			input:    []interface{}{"/admin"},
			want:     "http://app.bhojpur.net/admin",
		},
		{
			original: "http://app.bhojpur.net?q=keyword",
			input:    []interface{}{"admin"},
			want:     "http://app.bhojpur.net/admin?q=keyword",
		},
		{
			original: "http://app.bhojpur.net/?q=keyword",
			input:    []interface{}{"admin"},
			want:     "http://app.bhojpur.net/admin?q=keyword",
		},
		{
			original: "http://app.bhojpur.net/?q=keyword",
			input:    []interface{}{"admin/"},
			want:     "http://app.bhojpur.net/admin/?q=keyword",
		},
	}
	for _, c := range cases {
		// u, _ := url.Parse(c.original)
		// context := Context{Context: &bhojpur.Context{Request: &http.Request{URL: u}}}
		got, err := JoinURL(c.original, c.input...)
		if c.err != nil {
			if err == nil || err.Error() != c.err.Error() {
				t.Errorf("got error %s; want %s", err, c.err)
			}
		} else {
			if err != nil {
				t.Error(err)
			}
			if got != c.want {
				t.Errorf("context.JoinURL = %s; c.want %s", got, c.want)
			}
		}
	}
}

func TestSortFormKeys(t *testing.T) {
	keys := []string{"BhojpurResource.Category", "BhojpurResource.Addresses[2].Address1", "BhojpurResource.Addresses[1].Address1", "BhojpurResource.Addresses[11].Address1", "BhojpurResource.Addresses[0].Address1", "BhojpurResource.Code", "BhojpurResource.ColorVariations[0].Color", "BhojpurResource.ColorVariations[0].ID", "BhojpurResource.ColorVariations[0].SizeVariations[2].AvailableQuantity", "BhojpurResource.ColorVariations[0].SizeVariations[11].AvailableQuantity", "BhojpurResource.ColorVariations[0].SizeVariations[22].AvailableQuantity", "BhojpurResource.ColorVariations[0].SizeVariations[3].AvailableQuantity", "BhojpurResource.ColorVariations[0].SizeVariations[4].AvailableQuantity", "BhojpurResource.ColorVariations[1].SizeVariations[0].AvailableQuantity", "BhojpurResource.ColorVariations[1].SizeVariations[1].AvailableQuantity", "BhojpurResource.ColorVariations[1].ID", "BhojpurResource.ColorVariations[0].SizeVariations[1].ID", "BhojpurResource.ColorVariations[0].SizeVariations[4].ID", "BhojpurResource.ColorVariations[0].SizeVariations[3].ID", "BhojpurResource.Z[0]"}

	SortFormKeys(keys)

	orderedKeys := []string{"BhojpurResource.Addresses[0].Address1", "BhojpurResource.Addresses[1].Address1", "BhojpurResource.Addresses[2].Address1", "BhojpurResource.Addresses[11].Address1", "BhojpurResource.Category", "BhojpurResource.Code", "BhojpurResource.ColorVariations[0].Color", "BhojpurResource.ColorVariations[0].ID", "BhojpurResource.ColorVariations[0].SizeVariations[1].ID", "BhojpurResource.ColorVariations[0].SizeVariations[2].AvailableQuantity", "BhojpurResource.ColorVariations[0].SizeVariations[3].AvailableQuantity", "BhojpurResource.ColorVariations[0].SizeVariations[3].ID", "BhojpurResource.ColorVariations[0].SizeVariations[4].AvailableQuantity", "BhojpurResource.ColorVariations[0].SizeVariations[4].ID", "BhojpurResource.ColorVariations[0].SizeVariations[11].AvailableQuantity", "BhojpurResource.ColorVariations[0].SizeVariations[22].AvailableQuantity", "BhojpurResource.ColorVariations[1].ID", "BhojpurResource.ColorVariations[1].SizeVariations[0].AvailableQuantity", "BhojpurResource.ColorVariations[1].SizeVariations[1].AvailableQuantity", "BhojpurResource.Z[0]"}

	if fmt.Sprint(keys) != fmt.Sprint(orderedKeys) {
		t.Errorf("ordered form keys should be \n%v\n, but got\n%v", orderedKeys, keys)
	}
}

func TestSafeJoin(t *testing.T) {
	pth1, err := SafeJoin("hello", "world")
	if err != nil || pth1 != "hello/world" {
		t.Errorf("no error should happen")
	}

	// test possible vulnerability https://snyk.io/research/zip-slip-vulnerability#go
	pth2, err := SafeJoin("hello", "../world")
	if err == nil || pth2 != "" {
		t.Errorf("no error should happen")
	}
}
