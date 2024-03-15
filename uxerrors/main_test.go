/**
 * file: uxerrors/main_test.go
 * author: theo technicguy
 * license: apache-2.0
 */

package uxerrors

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateSeed(seed int64) int64 {
	if seed != 0 {
		return seed
	}
	return time.Now().UnixMilli()
}

func generateError(seed int64) (Error, int64) {
	s := generateSeed(seed)
	fkr := faker.NewWithSeed(rand.NewSource(s))
	return Error{
		Summary: fkr.Lorem().Sentence(10),
		Detail:  fkr.Lorem().Paragraph(3),
		Debug: Debug{
			Raw: fkr.Lorem().Text(64),
		},
	}, s
}

func generateErrors(count int, seed int64) (Errors, int64) {
	s := generateSeed(seed)

	errs := []Error{}
	for i := 0; i < count; i++ {
		e, _ := generateError(s)
		errs = append(errs, e)
	}

	return Errors{
		Status: faker.New().IntBetween(0, 599),
		Errors: errs,
	}, s
}

func TestNewError(t *testing.T) {
	seed := generateSeed(0)
	t.Logf("Using seed %x\n", seed)

	for i := 0; i < 10; i++ {
		expect, _ := generateError(seed)
		expect.Summary = ""
		expect.Detail = ""

		actual := New(fmt.Errorf(expect.Debug.Raw))
		assert.NotNil(t, actual, "a new error should not be nil")
		assert.Equal(t, expect, actual, "the expected and actual error structures should be the same")
	}
}

func TestPrintError(t *testing.T) {
	seed := generateSeed(0)
	t.Logf("Using seed %x\n", seed)

	for i := 0; i < 10; i++ {
		expect, _ := generateError(seed)

		actual := expect.Error()
		assert.NotNil(t, actual, "an error should not be nil")
		assert.NotEmpty(t, actual, "an error should not be empty")
		assert.Equal(t, expect.Error(), actual, "the expected and actual error string should be the same")
	}
}

func TestErrorToMapNoDebug(t *testing.T) {
	seed := generateSeed(0)
	t.Logf("Using seed %x\n", seed)

	for i := 0; i < 10; i++ {
		err, _ := generateError(seed)
		expect := map[string]any{
			"Summary": err.Summary,
			"Detail":  err.Detail,
		}

		actual := err.ToMap(false)
		assert.NotNil(t, actual, "an error map should not be nil")
		assert.NotEmpty(t, actual, "an error map should not be empty")
		assert.Equal(t, expect, actual, "the expected and actual error maps should be the same")
	}
}

func TestErrorToMapDebug(t *testing.T) {
	seed := generateSeed(0)
	t.Logf("Using seed %x\n", seed)

	for i := 0; i < 10; i++ {
		err, _ := generateError(seed)
		expect := map[string]any{
			"Summary": err.Summary,
			"Detail":  err.Detail,
			"Debug":   err.Debug,
		}

		actual := err.ToMap(true)
		assert.NotNil(t, actual, "an error map should not be nil")
		assert.NotEmpty(t, actual, "an error map should not be empty")
		assert.Equal(t, expect, actual, "the expected and actual error maps should be the same")
	}
}

func TestNewErrors(t *testing.T) {
	seed := generateSeed(0)
	t.Logf("Using seed %x\n", seed)

	for i := 0; i < 10; i++ {
		expect, _ := generateErrors(0, seed)

		actual := NewErrors(expect.Status)
		assert.NotNil(t, actual, "a new errors should not be nil")
		assert.Equal(t, expect, actual, "the expected and actual errors structure should be the same")
	}
}

func TestPrintErrors1(t *testing.T) {
	seed := generateSeed(0)
	t.Logf("Using seed %x\n", seed)

	for i := 0; i < 10; i++ {
		expect, _ := generateErrors(1, seed)
		expectErr := expect.Errors[0].Debug.Raw

		actual := NewErrors(expect.Status)
		assert.NotNil(t, actual, "a new errors should not be nil")

		actual.Errors = expect.Errors
		assert.Equal(t, expect, actual, "the expected and actual errors structure should be the same")
		assert.Equal(t, expectErr, actual.Error(), "the expected and actual error messages should be the same")
	}
}

func TestPrintErrors5(t *testing.T) {
	seed := generateSeed(0)
	t.Logf("Using seed %x\n", seed)

	for i := 0; i < 10; i++ {
		errs, _ := generateErrors(5, seed)
		expect := ""
		for _, v := range errs.Errors {
			expect += fmt.Sprintf("%s\n", v.Debug.Raw)
		}

		actual := NewErrors(errs.Status)
		assert.NotNil(t, actual, "a new errors should not be nil")

		actual.Errors = errs.Errors
		assert.Equal(t, errs, actual, "the expected and actual errors structure should be the same")
		assert.Equal(t, expect, actual.Error(), "the expected and actual error messages should be the same")
	}
}

func TestErrorsToMapNoDebug(t *testing.T) {
	seed := generateSeed(0)
	t.Logf("Using seed %x\n", seed)

	for i := 0; i < 10; i++ {
		errs, _ := generateErrors(5, seed)
		expect := map[string][]map[string]any{
			"Errors": {},
		}
		for _, v := range errs.Errors {
			expect["Errors"] = append(expect["Errors"], v.ToMap(false))
		}

		actual := NewErrors(errs.Status)
		assert.NotNil(t, actual, "a new errors should not be nil")

		actual.Errors = errs.Errors
		assert.Equal(t, errs, actual, "the expected and actual errors structure should be the same")
		assert.Equal(t, expect, actual.ToMap(false), "the expected and actual error maps should be the same")
	}
}

func TestErrorsToMapDebug(t *testing.T) {
	seed := generateSeed(0)
	t.Logf("Using seed %x\n", seed)

	for i := 0; i < 10; i++ {
		errs, _ := generateErrors(5, seed)
		expect := map[string][]map[string]any{
			"Errors": {},
		}
		for _, v := range errs.Errors {
			expect["Errors"] = append(expect["Errors"], v.ToMap(true))
		}

		actual := NewErrors(errs.Status)
		assert.NotNil(t, actual, "a new errors should not be nil")

		actual.Errors = errs.Errors
		assert.Equal(t, errs, actual, "the expected and actual errors structure should be the same")
		assert.Equal(t, expect, actual.ToMap(true), "the expected and actual error maps should be the same")
	}
}

func TestErrorsAppendOnEmpty(t *testing.T) {
	seed := generateSeed(0)
	t.Logf("Using seed %x\n", seed)

	for i := 0; i < 10; i++ {
		errs, _ := generateErrors(0, seed)
		require.NotNil(t, errs, "could not generate a new Errors")
		err, _ := generateError(seed)
		require.NotNil(t, err, "could not generate a new Error")

		expect := Errors{
			Status: errs.Status,
			Errors: []Error{err},
		}
		actual := errs.Append(err)

		assert.NotNil(t, actual, "the append return value should not be nil")
		assert.Equal(t, expect, actual, "the expected and return Errors should be the same")
	}
}

func TestErrorsAppendOnFilled(t *testing.T) {
	seed := generateSeed(0)
	t.Logf("Using seed %x\n", seed)

	for i := 0; i < 10; i++ {
		errs, _ := generateErrors(5, seed)
		require.NotNil(t, errs, "could not generate a new Errors")
		err, _ := generateError(seed)
		require.NotNil(t, err, "could not generate a new Error")

		expect := Errors{
			Status: errs.Status,
			Errors: append(errs.Errors, err),
		}
		actual := errs.Append(err)

		assert.NotNil(t, actual, "the append return value should not be nil")
		assert.Equal(t, expect, actual, "the expected and return Errors should be the same")
	}
}

func TestErrorsAppendNewOnEmpty(t *testing.T) {
	seed := generateSeed(0)
	t.Logf("Using seed %x\n", seed)

	for i := 0; i < 10; i++ {
		errs, _ := generateErrors(0, seed)
		require.NotNil(t, errs, "could not generate a new Errors")
		err, _ := generateError(seed)
		require.NotNil(t, err, "could not generate a new Error")
		msg := fmt.Errorf("%s", err.Debug.Raw)

		expect := Errors{
			Status: errs.Status,
			Errors: []Error{{Debug: Debug{Raw: err.Debug.Raw}}},
		}
		actual := errs.AppendNew(msg)

		assert.NotNil(t, actual, "the append return value should not be nil")
		assert.Equal(t, expect, actual, "the expected and return Errors should be the same")
	}
}

func TestErrorsAppendNewOnFilled(t *testing.T) {
	seed := generateSeed(0)
	t.Logf("Using seed %x\n", seed)

	for i := 0; i < 10; i++ {
		errs, _ := generateErrors(5, seed)
		require.NotNil(t, errs, "could not generate a new Errors")
		err, _ := generateError(seed)
		require.NotNil(t, err, "could not generate a new Error")
		msg := fmt.Errorf("%s", err.Debug.Raw)

		expect := Errors{
			Status: errs.Status,
			Errors: append(errs.Errors, err),
		}
		actual := errs.AppendNew(msg)

		assert.NotNil(t, actual, "the append return value should not be nil")
		assert.Equal(t, expect, actual, "the expected and return Errors should be the same")
	}
}
