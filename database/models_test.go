package database

import (
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCourseValidator(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.RegisterValidation("iscourse", isCourse, false)
	require.NoError(t, err, "could not register validator")

	validCourses := []string{
		"LINFO1000", "LEPL2020", "LBIR1210", "LTHEO4099", "LCOPS1509", "LINMA1702", "LPHYL9999",
	}

	for _, c := range validCourses {
		c = strings.ToLower(c)
		assert.NoErrorf(t, validate.Var(c, "iscourse"), "%s is a valid course\n", c)
	}

	invalidCourses := []string{
		"SINF11BA", "", "LINFO0000", "LTHECO666", "LIN001000",
	}

	for _, c := range invalidCourses {
		c = strings.ToLower(c)
		assert.Errorf(t, validate.Var(c, "iscourse"), "%s is an invalid course\n", c)
	}
}
