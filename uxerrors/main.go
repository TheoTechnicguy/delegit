/**
 * file: uxerrors/main.go
 * author: theo technciguy
 * license: apache-2.0
 *
 * The uxerrors package provides a more user-friendly
 * error handling experience for the application.
 */

// UXErrors are errors that are more user-friendly,
// as much for non-technical users as for application
// developers.
package uxerrors

import (
	"fmt"
)

// The ErrorsX Error structure is a more specific implementation of
// error, designed for an improved user and developer experience.
//
// Error implements the error interface.
type Error struct {
	// Summary contains a short, user friendly message that
	// understandably describes the error that occurred.
	Summary string `json:"Summary"`

	// Detail contains a more verbose message that describes
	// why the error occurred and how to resolve it.
	Detail string `json:"Detail"`

	// Debug contains additional debug information. It should
	// generally not be displayed to users, and is intended
	// for more technically experienced users.
	Debug Debug `json:"Debug"`
}

// The Debug structure is designed to assist with debugging
// applications, in development as well as in production
// environments. It contains the raw (original) error that
// created the Error, and additional information that may
// help finding the source of the error.
type Debug struct {
	// Raw is the raw, unprocessed error that generated the
	// Error
	Raw string `json:"Raw"`
}

// New creates a new UXError from a generic error. The error
// may not be nil. The Summary and Detail fields are not set,
// and it is up to the caller to set them correctly.
func New(err error) Error {
	return Error{
		Debug: Debug{
			Raw: err.Error(),
		},
	}
}

// Error returns the raw error string. It implements the error
// interface.
func (e Error) Error() string {
	return e.Debug.Raw
}

// ToMap returns a map representation of the Error.
// Depending on the value for debug, the Debug field
// may be included in the map. Including it is useful
// for debugging. You may also want to omit debug
// information when returning errors to users.
func (e Error) ToMap(debug bool) (m map[string]any) {
	m = map[string]any{
		"Summary": e.Summary,
		"Detail":  e.Detail,
	}

	if debug {
		m["Debug"] = e.Debug
	}

	return
}

// The Errors structure is a collection of Error[s]
// containing a status code for the HTTP server.
// It collects all errors in an array so that they can
// all be returned to the user in one go, improving UX.
type Errors struct {
	// Status is the HTTP status code that should be
	// returned to the user.
	Status int

	// Errors is an array of Error[s] collecting all
	// errors that occurred.
	// The entire array is returned to the user.
	Errors []Error `json:"Errors"`
}

// NewErrors creates a new Errors structure with the
// given status code. The Errors array is empty. You
// can use the Append method to add errors.
func NewErrors(status int) Errors {
	es := Errors{Status: status, Errors: []Error{}}
	return es
}

// Error returns an aggregation of all errors in the
// array. It implements the error interface.
func (es Errors) Error() string {
	if len(es.Errors) == 1 {
		return es.Errors[0].Debug.Raw
	}

	s := ""
	for _, v := range es.Errors {
		s += fmt.Sprintf("%s\n", v.Error())
	}

	return s
}

// ToMap returns a map representation of the Errors.
// Depending on the value for debug, the Debug field
// may be included in the map. Including it is useful
// for debugging. You may also want to omit debug
// information when returning errors to users.
func (es Errors) ToMap(debug bool) map[string][]map[string]any {
	m := map[string][]map[string]any{}
	m["Errors"] = []map[string]any{}

	for _, v := range es.Errors {
		m["Errors"] = append(m["Errors"], v.ToMap(debug))
	}

	return m
}

// Append adds an error to the Errors array and returns
// the new Errors structure.
func (es Errors) Append(err Error) Errors {
	es.Errors = append(es.Errors, err)
	return es
}

// AppendNew adds a new error to the Errors array and
// returns the new Errors structure. The error is created
// from the given error.
// In effect, it comines the Error.New and Errors.Append
// methods.
func (es Errors) AppendNew(err error) Errors {
	return es.Append(New(err))
}
