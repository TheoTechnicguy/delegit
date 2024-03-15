/**
 * file: validators/utils.go
 * author: theo technciguy
 * license: apache-2.0
 *
 * The utils file contains package-widely used
 * utility functions for the validators
 * package.
 */

package validators

import (
	"fmt"
	"reflect"

	"git.licolas.net/delegit/delegit/uxerrors"
	"github.com/go-playground/validator/v10"
)

// requiredMissingError sets the summary and detail fields of the
// error to indicate that a required field is missing.
func requiredMissingError(xerr *uxerrors.Error, err validator.FieldError) {
	xerr.Summary = fmt.Sprintf("The %s field is missing", err.Field())
	xerr.Detail = fmt.Sprintf("The %s field is a required field, however it is empty. Fill the field correctly and try again.", err.Field())
}

// genericError sets the summary and detail fields of the error to
// indicate that an unspecified error occurred while validating a
// field.
func genericError(xerr *uxerrors.Error, err validator.FieldError) {
	xerr.Summary = fmt.Sprintf("There was an error while validating the %s field", err.Field())
	xerr.Detail = fmt.Sprintf("While validating the %s field, an unspecified error occurred. Check the field for correctness and try again.", err.Field())
}

// minError sets the summary and detail fields of the error to
// indicate that a string field is too short, or a number field
// is too small.
func minError(xerr *uxerrors.Error, err validator.FieldError) {
	var summary, detail string

	switch err.Type() {
	case reflect.TypeFor[string]():
		summary = "is too short"
		detail = fmt.Sprintf("It should be at least %s long, but was %d. Elaborate and try again.", err.Param(), len(err.Value().(string)))
	case reflect.TypeFor[int](), reflect.TypeFor[uint]():
		summary = "is too small"
		detail = fmt.Sprintf("It should be at least %s, but was %d. Increase the value and try again.", err.Param(), err.Value())
	}

	xerr.Summary = fmt.Sprintf("The %s field %s", err.Field(), summary)
	xerr.Detail = fmt.Sprintf("The %s field %s. %s", err.Field(), summary, detail)
}

// maxError sets the summary and detail fields of the error to
// indicate that a string field is too long, or a number field
// is too large.
func maxError(xerr *uxerrors.Error, err validator.FieldError) {
	var summary, detail string

	switch err.Type() {
	case reflect.TypeFor[string]():
		summary = "is too long"
		detail = fmt.Sprintf("It should be at most %s long, but was %d. Shorten and try again.", err.Param(), len(err.Value().(string)))
	case reflect.TypeFor[int](), reflect.TypeFor[uint]():
		summary = "is too high"
		detail = fmt.Sprintf("It should be at most %s, but was %d. Decrease the value and try again.", err.Param(), err.Value())
	}

	xerr.Summary = fmt.Sprintf("The %s field %s", err.Field(), summary)
	xerr.Detail = fmt.Sprintf("The %s field %s. %s", err.Field(), summary, detail)
}
