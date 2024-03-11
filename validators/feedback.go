/**
 * file: validators/course.go
 * author: theo technciguy
 * license: apache-2.0
 *
 * The course validator validates the course field of the
 * feedback form.
 */

// Package validators provides validation functions for the
// feedback form.
// It also includes custom validators for specific fields.
package validators

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"git.licolas.net/delegit/delegit/models"
	"git.licolas.net/delegit/delegit/uxerrors"
	"github.com/go-playground/validator/v10"
)

// IsCourse validates that a filed fills the UCLouvain course
// scheme. The course scheme is as follows.
//
//	course := "L" faculty code
//	faculty := letter letter letter letter?
//	code := digit digit digit digit
//
// Letters are alpha ascii letters, generally uppercase. Lowercase
// formatting should be accepted however.
// Codes are at least 1000 and at most 9999.
func IsCourse(fl validator.FieldLevel) bool {
	course := strings.ToLower(fl.Field().String())

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Var(course, "alphanum,min=8,max=9,startswith=l"); err != nil {
		return false
	}

	code, err := strconv.Atoi(course[len(course)-4:])
	if err != nil {
		return false
	}
	if code < 1000 {
		return false
	}

	faculty := course[1 : len(course)-4]
	if err := validate.Var(faculty, "alpha,min=3,max=4"); err != nil {
		return false
	}

	return true
}

func IsWhitespace(fl validator.FieldLevel) bool {
	text := fl.Field().String()
	for _, v := range text {
		switch v {
		case ' ', '\t', '\n', '\v', '\f', '\r':
			continue
		default:
			return false
		}
	}
	return true
}

func IsPunctuation(fl validator.FieldLevel) bool {
	text := fl.Field().String()
	for _, v := range text {
		switch v {
		case '.', ',', ';', ':', '!', '?', '(', ')', '[', ']', '{', '}', '<', '>', '"', '\'', '/', '\\', '|', '@', '#', '$', '%', '^', '&', '*', '-', '_', '=', '+', '~', '`':
			continue
		default:
			return false
		}
	}
	return true
}

func IsAsciiNumUnicodeText(fl validator.FieldLevel) bool {
	feedback := strings.ToLower(fl.Field().String())
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("whitespace", IsWhitespace)
	validate.RegisterValidation("punctuation", IsPunctuation)
	for _, v := range feedback {
		if err := validate.Var(string(v), "alphanumunicode|whitespace|punctuation"); err != nil {
			return false
		}
	}
	return true
}

// ValidateFeedback validates the feedback structure. It returns an
// UXErrors containing all the errors that occurred during validation
// or nil if no errors occurred.
func ValidateFeedback(f *models.Feedback) error {
	v := validator.New()
	v.RegisterValidation("iscourse", IsCourse, false)
	v.RegisterValidation("alphanumunicodetext", IsAsciiNumUnicodeText, false)
	err := v.Struct(f)

	if err == nil {
		return nil
	}

	vErr := err.(validator.ValidationErrors)
	errs := uxerrors.Errors{Status: http.StatusBadRequest}
	for _, ve := range vErr {
		xerr := uxerrors.New(err)

		switch ve.Tag() {
		case "required":
			requiredMissingError(&xerr, ve)
		case "alphanumunicodetext":
			xerr.Summary = fmt.Sprintf("The %s field contains forbidden characters", ve.Field())
			xerr.Detail = fmt.Sprintf("The %s field contains forbidden characters. Only letter, numbers and special characters are allowed. Remove all others and try again", ve.Field())
		case "min", "ge", "gt":
			minError(&xerr, ve)
		case "max", "le", "lt":
			maxError(&xerr, ve)
		case "iscourse":
			xerr.Summary = "The course does not look like a valid course"
			xerr.Detail = fmt.Sprintf("The course you entered (%q) does not look like a valid course code. Check the code and try again.", ve.Value())
		default:
			genericError(&xerr, ve)
		}

		errs.Errors = append(errs.Errors, xerr)
	}

	return errs
}
