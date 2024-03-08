package validators

import (
	"fmt"
	"net/http"
	"reflect"
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

func ValidateFeedback(f *models.Feedback) error {
	v := validator.New()
	v.RegisterValidation("iscourse", IsCourse, false)
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
		case "alphanumunicode":
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

func requiredMissingError(xerr *uxerrors.Error, err validator.FieldError) {
	xerr.Summary = fmt.Sprintf("The %s field is missing", err.Field())
	xerr.Detail = fmt.Sprintf("The %s field is a required field, however it is empty. Fill the field correctly and try again.", err.Field())
}

func genericError(xerr *uxerrors.Error, err validator.FieldError) {
	xerr.Summary = fmt.Sprintf("There was an error while validating the %s field", err.Field())
	xerr.Detail = fmt.Sprintf("While validating the %s field, an unspecified error occurred. Check the field for correctness and try again.", err.Field())
}

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

func maxError(xerr *uxerrors.Error, err validator.FieldError) {
	var summary, detail string

	switch err.Type() {
	case reflect.TypeFor[string]():
		summary = "is too long"
		detail = fmt.Sprintf("It should be at most %s long, but was %d. Shorten and try again.", err.Param(), len(err.Value().(string)))
	case reflect.TypeFor[int](), reflect.TypeFor[uint]():
		summary = "is too small"
		detail = fmt.Sprintf("It should be at most %s, but was %d. Decrease the value and try again.", err.Param(), err.Value())
	}

	xerr.Summary = fmt.Sprintf("The %s field %s", err.Field(), summary)
	xerr.Detail = fmt.Sprintf("The %s field %s. %s", err.Field(), summary, detail)
}
