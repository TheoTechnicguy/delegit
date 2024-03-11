package validators

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"git.licolas.net/delegit/delegit/models"
	"git.licolas.net/delegit/delegit/uxerrors"
	"github.com/go-playground/validator/v10"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// generateCourse generates a random, valid course
// code, following the UCLouvain course scheme.
func generateCourse(fkr faker.Faker) string {
	facultyLength := fkr.IntBetween(3, 4)
	var course strings.Builder

	course.WriteString("L")

	for i := 0; i < facultyLength; i++ {
		course.WriteString(fkr.Letter())
	}

	course.WriteString(strconv.Itoa(fkr.IntBetween(1000, 9999)))

	return course.String()
}

// generateFeedback is a helper function for tests.
// generating feedback using Faker, and returning a
// feedback. It takes a faker.Faker to generate the
// feedback.  The mutator function allows to mutate
// the feedback structure in place, after it has been
// filled by Faker.
//
// The function returns the feedback.
func generateFeedback(fkr faker.Faker, mutator func(*models.Feedback, faker.Faker)) *models.Feedback {
	if mutator == nil {
		mutator = func(f *models.Feedback, fkr faker.Faker) {}
	}

	newFeedback := &models.Feedback{
		ID:        fkr.UInt(),
		Course:    generateCourse(fkr),
		Feedback:  fkr.Lorem().Paragraph(3),
		Upvotes:   fkr.UIntBetween(0, 2000),
		Downvotes: fkr.UIntBetween(0, 2000),
	}
	mutator(newFeedback, fkr)

	return newFeedback
}

// generateFeedbackArray is a helper function for tests.
// It generates an array of n feedbacks using Faker, and
// returns the array. It takes a faker.Faker to generate
// the feedbacks. The mutator function allows to mutate
// the feedback structure in place, after it has been
// filled by Faker.
func generateFeedbackArray(fkr faker.Faker, mutator func(*models.Feedback, faker.Faker), n uint) (f []*models.Feedback) {
	var i uint = 0
	for ; i < n; i++ {
		f = append(f, generateFeedback(fkr, mutator))
	}

	return
}

// TestCourseValidator tests the IsCourse validator
// on a set of predefined valid and invalid courses.
func TestCourseValidator(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.RegisterValidation("iscourse", IsCourse, false)
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

// TestCourseValidatorFaker tests the IsCourse validator
// on a set of randomly generated courses.
func TestCourseValidatorFaker(t *testing.T) {
	seed := time.Now().UnixMilli()
	t.Logf("Current seed: %d\n", seed)

	fkr := faker.NewWithSeed(rand.NewSource(seed))
	val := validator.New(validator.WithRequiredStructEnabled())
	val.RegisterValidation("iscourse", IsCourse, false)

	for i := 0; i < 10; i++ {
		course := generateCourse(fkr)
		err := val.Var(course, "iscourse")
		assert.NoError(t, err, "the course generated should be valid")
	}
}

func TestValidateFeedbackValid(t *testing.T) {
	seed := time.Now().UnixMilli()
	t.Logf("Current seed: %d\n", seed)

	fkr := faker.NewWithSeed(rand.NewSource(seed))
	for i := 0; i < 10; i++ {
		feedback := generateFeedback(fkr, nil)
		t.Logf("feedback: %q\n", feedback.Feedback)
		err := ValidateFeedback(feedback)
		assert.NoError(t, err, "the feedback generated is valid and should pass validation")
	}
}

func TestValidateFeedbackMissingCourse(t *testing.T) {
	seed := time.Now().UnixMilli()
	t.Logf("Current seed: %d\n", seed)

	fkr := faker.NewWithSeed(rand.NewSource(seed))
	mutator := func(f *models.Feedback, fkr faker.Faker) {
		f.Course = ""
	}
	for i := 0; i < 10; i++ {
		feedback := generateFeedback(fkr, mutator)
		err := ValidateFeedback(feedback)
		assert.Error(t, err, "the feedback generated is missing a course and should not pass validation")
		require.IsType(t, uxerrors.Errors{}, err, "the error returned should be of type uxerrors.Errors")
		xerr := err.(uxerrors.Errors)
		assert.Len(t, xerr.Errors, 1, "the error returned should contain only one error")
		assert.Equal(t, "The Course field is missing", xerr.Errors[0].Summary, "the error returned should indicate that the course field is missing")
		assert.Equal(t, "The Course field is a required field, however it is empty. Fill the field correctly and try again.", xerr.Errors[0].Detail, "the error returned should indicate that the course field is missing")
	}
}

func TestValidateFeedbackMissingFeedback(t *testing.T) {
	seed := time.Now().UnixMilli()
	t.Logf("Current seed: %d\n", seed)

	fkr := faker.NewWithSeed(rand.NewSource(seed))
	mutator := func(f *models.Feedback, fkr faker.Faker) {
		f.Feedback = ""
	}
	for i := 0; i < 10; i++ {
		feedback := generateFeedback(fkr, mutator)
		err := ValidateFeedback(feedback)
		assert.Error(t, err, "the feedback generated is missing a feedback and should not pass validation")
		require.IsType(t, uxerrors.Errors{}, err, "the error returned should be of type uxerrors.Errors")
		xerr := err.(uxerrors.Errors)
		assert.Len(t, xerr.Errors, 1, "the error returned should contain only one error")
		assert.Equal(t, "The Feedback field is missing", xerr.Errors[0].Summary, "the error returned should indicate that the feedback field is missing")
		assert.Equal(t, "The Feedback field is a required field, however it is empty. Fill the field correctly and try again.", xerr.Errors[0].Detail, "the error returned should indicate that the feedback field is missing")
	}
}

func TestValidateFeedbackInvalidCourse(t *testing.T) {
	seed := time.Now().UnixMilli()
	t.Logf("Current seed: %d\n", seed)

	fkr := faker.NewWithSeed(rand.NewSource(seed))
	mutator := func(f *models.Feedback, fkr faker.Faker) {
		f.Course = "SINF11BA"
	}
	for i := 0; i < 10; i++ {
		feedback := generateFeedback(fkr, mutator)
		err := ValidateFeedback(feedback)
		assert.Error(t, err, "the feedback generated contains an invalid course and should not pass validation")
		require.IsType(t, uxerrors.Errors{}, err, "the error returned should be of type uxerrors.Errors")
		xerr := err.(uxerrors.Errors)
		assert.Len(t, xerr.Errors, 1, "the error returned should contain only one error")
		assert.Equal(t, "The course does not look like a valid course", xerr.Errors[0].Summary, "the error returned should indicate that the course field is invalid")
		assert.Equal(t, "The course you entered (\"SINF11BA\") does not look like a valid course code. Check the code and try again.", xerr.Errors[0].Detail, "the error returned should indicate that the course field is invalid")
	}
}

func TestValidateFeedbackInvalidFeedback(t *testing.T) {
	seed := time.Now().UnixMilli()
	t.Logf("Current seed: %d\n", seed)

	fkr := faker.NewWithSeed(rand.NewSource(seed))
	mutator := func(f *models.Feedback, fkr faker.Faker) {
		f.Feedback = "大家好！我叫小龙。我很喜欢沙拉。laksjflkjsakjflksajkfjlksajflkjdsalkfjlkajflkjdsalkfjlkdsajf dsljfajslk"
	}
	for i := 0; i < 10; i++ {
		feedback := generateFeedback(fkr, mutator)
		err := ValidateFeedback(feedback)
		assert.Error(t, err, "the feedback generated contains an invalid feedback and should not pass validation")
		require.IsType(t, uxerrors.Errors{}, err, "the error returned should be of type uxerrors.Errors")
		xerr := err.(uxerrors.Errors)
		assert.Len(t, xerr.Errors, 1, "the error returned should contain only one error")
		assert.Equal(t, "The Feedback field contains forbidden characters", xerr.Errors[0].Summary, "the error returned should indicate that the feedback field contains forbidden characters")
		assert.Equal(t, "The Feedback field contains forbidden characters. Only letter, numbers and special characters are allowed. Remove all others and try again", xerr.Errors[0].Detail, "the error returned should indicate that the feedback field contains forbidden characters")
	}
}

func TestValidateFeedbackTooShortFeedback(t *testing.T) {
	seed := time.Now().UnixMilli()
	t.Logf("Current seed: %d\n", seed)

	fkr := faker.NewWithSeed(rand.NewSource(seed))
	mutator := func(f *models.Feedback, fkr faker.Faker) {
		f.Feedback = "hi"
	}
	for i := 0; i < 10; i++ {
		feedback := generateFeedback(fkr, mutator)
		err := ValidateFeedback(feedback)
		assert.Error(t, err, "the feedback generated contains a too short feedback and should not pass validation")
		require.IsType(t, uxerrors.Errors{}, err, "the error returned should be of type uxerrors.Errors")
		xerr := err.(uxerrors.Errors)
		assert.Len(t, xerr.Errors, 1, "the error returned should contain only one error")
		assert.Equal(t, "The Feedback field is too short", xerr.Errors[0].Summary, "the error returned should indicate that the feedback field is too short")
		assert.Equal(t, "The Feedback field is too short. It should be at least 25 long, but was 2. Elaborate and try again.", xerr.Errors[0].Detail, "the error returned should indicate that the feedback field is too short")
	}
}

func TestValidateFeedbackTooLongFeedback(t *testing.T) {
	seed := time.Now().UnixMilli()
	t.Logf("Current seed: %d\n", seed)

	fkr := faker.NewWithSeed(rand.NewSource(seed))
	mutator := func(f *models.Feedback, fkr faker.Faker) {
		f.Feedback = fkr.Lorem().Paragraph(30)
	}
	for i := 0; i < 10; i++ {
		feedback := generateFeedback(fkr, mutator)
		err := ValidateFeedback(feedback)
		assert.Error(t, err, "the feedback generated contains a too long feedback and should not pass validation")
		require.IsType(t, uxerrors.Errors{}, err, "the error returned should be of type uxerrors.Errors")
		xerr := err.(uxerrors.Errors)
		assert.Len(t, xerr.Errors, 1, "the error returned should contain only one error")
		assert.Equal(t, "The Feedback field is too long", xerr.Errors[0].Summary, "the error returned should indicate that the feedback field is too long")
		assert.Equal(t, fmt.Sprintf("The Feedback field is too long. It should be at most 2000 long, but was %d. Shorten and try again.", len(feedback.Feedback)), xerr.Errors[0].Detail, "the error returned should indicate that the feedback field is too long")

	}
}

func TestValidateFeedbackTooBigUpvotes(t *testing.T) {
	seed := time.Now().UnixMilli()
	t.Logf("Current seed: %d\n", seed)

	fkr := faker.NewWithSeed(rand.NewSource(seed))
	mutator := func(f *models.Feedback, fkr faker.Faker) {
		f.Upvotes = 2001
	}
	for i := 0; i < 10; i++ {
		feedback := generateFeedback(fkr, mutator)
		err := ValidateFeedback(feedback)
		assert.Error(t, err, "the feedback generated contains an invalid upvotes and should not pass validation")
		require.IsType(t, uxerrors.Errors{}, err, "the error returned should be of type uxerrors.Errors")
		xerr := err.(uxerrors.Errors)
		assert.Len(t, xerr.Errors, 1, "the error returned should contain only one error")
		assert.Equal(t, "The Upvotes field is too high", xerr.Errors[0].Summary, "the error returned should indicate that the upvotes field is too high")
		assert.Equal(t, "The Upvotes field is too high. It should be at most 2000, but was 2001. Decrease the value and try again.", xerr.Errors[0].Detail, "the error returned should indicate that the upvotes field is too high")
	}
}

func TestValidateFeedbackTooBigDownvotes(t *testing.T) {
	seed := time.Now().UnixMilli()
	t.Logf("Current seed: %d\n", seed)

	fkr := faker.NewWithSeed(rand.NewSource(seed))
	mutator := func(f *models.Feedback, fkr faker.Faker) {
		f.Downvotes = 2001
	}
	for i := 0; i < 10; i++ {
		feedback := generateFeedback(fkr, mutator)
		err := ValidateFeedback(feedback)
		assert.Error(t, err, "the feedback generated contains an invalid downvotes and should not pass validation")
		require.IsType(t, uxerrors.Errors{}, err, "the error returned should be of type uxerrors.Errors")
		xerr := err.(uxerrors.Errors)
		assert.Len(t, xerr.Errors, 1, "the error returned should contain only one error")
		assert.Equal(t, "The Downvotes field is too high", xerr.Errors[0].Summary, "the error returned should indicate that the downvotes field is too high")
		assert.Equal(t, "The Downvotes field is too high. It should be at most 2000, but was 2001. Decrease the value and try again.", xerr.Errors[0].Detail, "the error returned should indicate that the downvotes field is too high")
	}
}
