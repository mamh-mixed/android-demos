package util

import (
	"bytes"
	"net/url"
	"sort"

	"github.com/google/go-querystring/query"
)

// Query implements encoding of struct into URL query parameters without escape
//
// Query expects to be passed a struct, and traverses it recursively using the
// following encoding rules.
//
// Each exported struct field is encoded as a URL parameter unless
//
//	- the field's tag is "-", or
//	- the field is empty and its tag specifies the "omitempty" option
//
// The empty values are false, 0, any nil pointer or interface value, any array
// slice, map, or string of length zero, and any time.Time that returns true
// for IsZero().
//
// The URL parameter name defaults to the struct field name but can be
// specified in the struct field's tag value.  The "url" key in the struct
// field's tag value is the key name, followed by an optional comma and
// options.  For example:
//
// 	// Field is ignored by this package.
// 	Field int `url:"-"`
//
// 	// Field appears as URL parameter "myName".
// 	Field int `url:"myName"`
//
// 	// Field appears as URL parameter "myName" and the field is omitted if
// 	// its value is empty
// 	Field int `url:"myName,omitempty"`
//
// 	// Field appears as URL parameter "Field" (the default), but the field
// 	// is skipped if empty.  Note the leading comma.
// 	Field int `url:",omitempty"`
//
// For encoding individual field values, the following type-dependent rules
// apply:
//
// Boolean values default to encoding as the strings "true" or "false".
// Including the "int" option signals that the field should be encoded as the
// strings "1" or "0".
//
// time.Time values default to encoding as RFC3339 timestamps.  Including the
// "unix" option signals that the field should be encoded as a Unix time (see
// time.Unix())
//
// Slice and Array values default to encoding as multiple URL values of the
// same name.  Including the "comma" option signals that the field should be
// encoded as a single comma-delimited value.  Including the "space" option
// similarly encodes the value as a single space-delimited string. Including
// the "brackets" option signals that the multiple URL values should have "[]"
// appended to the value name.
//
// Anonymous struct fields are usually encoded as if their inner exported
// fields were fields in the outer struct, subject to the standard Go
// visibility rules.  An anonymous struct field with a name given in its URL
// tag is treated as having that name, rather than being anonymous.
//
// Non-nil pointer values are encoded as the value pointed to.
//
// Nested structs are encoded including parent fields in value names for
// scoping. e.g:
//
// 	"user[name]=acme&user[addr][postcode]=1234&user[addr][city]=SFO"
//
// All other values are encoded using their default string representation.
//
// Multiple fields that encode to the same URL parameter name will be included
// as multiple URL values of the same name.
func Query(s interface{}, excludes ...string) (buf bytes.Buffer, err error) {
	if s == nil {
		return
	}

	v, err := query.Values(s)
	if err != nil {
		return buf, err
	}

	return QueryValues(v), nil
}

// QueryValues implements encoding of values into URL query parameters without escape
func QueryValues(v url.Values, excludes ...string) (buf bytes.Buffer) {
	if v == nil {
		return
	}

	keys := make([]string, 0, len(v))
	for k := range v {
		if StringInSlice(k, excludes) {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		vs := v[k]
		prefix := k + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	return buf

}
