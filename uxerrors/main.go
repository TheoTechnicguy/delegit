/**
 * file: uxerrors/main.go
 * author: theo technciguy
 * license: apache-2.0
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

func New(err error) Error {
	return Error{
		Debug: Debug{
			Raw: err.Error(),
		},
	}
}

func (e Error) Error() string {
	return e.Debug.Raw
}

func (e Error) ToMap(debug bool) map[string]any {
	m := map[string]any{}
	m["Summary"] = e.Summary
	m["Detail"] = e.Detail

	if debug {
		m["Debug"] = e.Debug
	}

	return m
}

type Errors struct {
	Status int
	Errors []Error `json:"Errors"`
}

func NewErrors(status int) Errors {
	es := Errors{Status: status}
	return es
}

func (es Errors) Error() string {
	s := ""
	for _, v := range es.Errors {
		s += fmt.Sprintf("%s\n", v.Error())
	}

	return s
}

func (es Errors) ToMap(debug bool) map[string][]map[string]any {
	m := map[string][]map[string]any{}
	m["Errors"] = []map[string]any{}

	for _, v := range es.Errors {
		m["Errors"] = append(m["Errors"], v.ToMap(debug))
	}

	return m
}

func (es Errors) Append(err Error) Errors {
	es.Errors = append(es.Errors, err)
	return es
}

func (es Errors) AppendNew(err error) Errors {
	return es.Append(New(err))
}
