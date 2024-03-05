package validators

import (
	"strconv"
	"strings"

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
