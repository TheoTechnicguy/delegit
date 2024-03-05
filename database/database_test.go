/**
 * file: database/database_test.go
 * author: theo technicguy
 * license: apache-2.0
 *
 * This file provides unit test cases for
 * the data persistence.
 *
 * A general note on the queries.
 * We will be assuming that GORM (and the community)
 * is doing it's job and testing their backend.
 * These tests are not here to test the queries
 * executed by GORM, but the arguments that are passed
 * to these queries.
 */

package database

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func createMockDatabase(t *testing.T) (*Database, func(), sqlmock.Sqlmock, *sqlmock.Rows) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "could not create database mock")

	dialect := postgres.New(postgres.Config{
		Conn: mockDB,
	})
	db, err := NewDatabaseFromDialector(dialect, &gorm.Config{})
	require.NoError(t, err, "could not create database from dialect")

	schema := sqlmock.NewRows([]string{"id", "course", "feedback", "upvotes", "downvotes"})

	closer := func() {
		mockDB.Close()
	}

	return db, closer, mock, schema
}

func feedbackToCSV(f ...*Feedback) (s string) {
	fs := []string{}
	for _, v := range f {
		fs = append(
			fs,
			fmt.Sprintf(
				"%d,%s,%s,%d,%d",
				v.ID,
				v.Course,
				v.Feedback,
				v.Upvotes,
				v.Downvotes,
			),
		)
	}
	return strings.Join(fs, "\n")
}

// generateFeedback is a helper function for tests.
// generating feedback using Faker, and returning a
// slice of `n` feedback. It takes a seed, which can
// be `0`. In that case, the seed will be generated.
// The mutator function allows to mutate the feedback
// structure in place, after it has been filled by
// Faker.
//
// The function returns the slice of feedback and the
// actual seed used.
func generateFeedback(n uint, seed int64, mutator func(*Feedback, faker.Faker)) (f []*Feedback, s int64) {
	if seed == 0 {
		s = time.Now().UnixMilli()
	} else {
		s = seed
	}

	if mutator == nil {
		mutator = func(f *Feedback, fkr faker.Faker) {}
	}

	fkr := faker.NewWithSeed(rand.NewSource(s))
	var i uint
	for i = 0; i < n; i++ {
		newFeedback := &Feedback{}
		newFeedback.ID = fkr.UInt()
		newFeedback.Course = fkr.RandomStringWithLength(10)
		newFeedback.Feedback = fkr.Lorem().Paragraph(3)
		newFeedback.Upvotes = fkr.UIntBetween(0, 2000)
		newFeedback.Downvotes = fkr.UIntBetween(0, 2000)
		mutator(newFeedback, fkr)
		f = append(f, newFeedback)
	}

	return
}

// TestGetAllFeedbackWhenEmpty is a unit test that tests the
// return value of GetAllFeedback when the database is empty.
// There are two cases a database can be empty:
//  1. It is newly created
//  2. Everything got deleted
//
// Either way, no feedback should be returned.
// Here, we test the first case.
func TestGetAllFeedbackDatabaseEmpty(t *testing.T) {
	db, closer, mock, schema := createMockDatabase(t)
	defer closer()

	mock.
		ExpectQuery("^SELECT .+ FROM [`\"']feedbacks[`\"'].*$").
		WithoutArgs().
		WillReturnRows(schema)

	fb, err := db.GetAllFeedback()
	assert.NoError(t, err, "fetching all feedback returned an error")
	assert.ElementsMatch(t, []*Feedback{}, fb, "new database fetch should not return any feedback")
}

// TestGetAllFeedback is a unit test that tests the return
// value of GetAllFeedback when the database has some
// feedback stored.
func TestGetAllFeedback(t *testing.T) {
	db, closer, mock, schema := createMockDatabase(t)
	defer closer()

	expectFeedback, seed := generateFeedback(10, 0, nil)
	t.Logf("seed: %x\n", seed)

	schema.FromCSVString(feedbackToCSV(expectFeedback...))
	mock.
		ExpectQuery("^SELECT .+ FROM [`\"']feedbacks[`\"'].*$").
		WithoutArgs().
		WillReturnRows(schema)

	fb, err := db.GetAllFeedback()
	assert.NoError(t, err, "fetching all feedback returned an error")
	assert.ElementsMatch(t, expectFeedback, fb, "feedback elements should match")
}

// TestGetFeedback is a unit test that tests the return value
// of the by id requested feedback.
// GetFeedback is expected to return the requested feedback,
// matching is done by ID, and no error whe the feedback exists,
// and no feedback + error when it is not found.
func TestGetFeedback(t *testing.T) {
	db, closer, mock, schema := createMockDatabase(t)
	defer closer()

	expectFeedback, seed := generateFeedback(10, 0, nil)
	t.Logf("seed: %x\n", seed)

	for _, expected := range expectFeedback {
		expectSQL := schema
		expectSQL.FromCSVString(feedbackToCSV(expected))
		mock.
			ExpectQuery("^SELECT .+ FROM [`\"']feedbacks[`\"'] WHERE [`\"']feedbacks[`\"']\\.[`\"']id[`\"']\\W*=.*$").
			WithArgs(expected.ID, 1).
			WillReturnRows(expectSQL)

		actual, err := db.GetFeedback(expected.ID)
		require.NoError(t, err, "fetching feedback failed")
		assert.Equal(t, expected, actual, "feedback are different but are supposed to be the same")
	}
}

// TestGetFeedbackNotExists is a unit test that tests the return
// value of the feedback requested by the passed id. More
// specifically, it is supposed to return an error.
// GetFeedback is expected to return the requested feedback,
// matching is done by ID, and no error whe the feedback exists,
// and no feedback + error when it is not found.
func TestGetFeedbackNotExists(t *testing.T) {
	db, closer, mock, _ := createMockDatabase(t)
	defer closer()

	mutator := func(f *Feedback, fkr faker.Faker) {
		f.ID = fkr.UIntBetween(0, 0xffff)
	}
	_, seed := generateFeedback(10, 0, mutator)
	t.Logf("seed: %x\n", seed)

	fkr := faker.NewWithSeed(rand.NewSource(seed))
	for i := 0; i < 10; i++ {
		id := fkr.UIntBetween(0x10000, 0xffffff)

		mock.
			ExpectQuery("^SELECT .+ FROM [`\"']feedbacks[`\"'] WHERE [`\"']feedbacks[`\"']\\.[`\"']id[`\"']\\W*=.*$").
			WithArgs(id, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		feedback, err := db.GetFeedback(id)
		assert.Error(t, err, "an error was expected")
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound, "incorrect error returned")
		assert.Nil(t, feedback, "there should be no feedback returned")
	}
}

// TestAddFeedback is a unit test that tests the return
// value and the stored data in the database using randomly
// created feedback (by faker).
// AddFeedback is expected to set the correct ID, and reset
// the Up- and Downvotes, add the sanitized data to the database
// and return this structure.
func TestAddFeedback(t *testing.T) {
	db, closer, mock, schema := createMockDatabase(t)
	defer closer()

	mutator := func(f *Feedback, fkr faker.Faker) {
		f.ID = fkr.UIntBetween(0, 0xffff)
	}
	expectedFeedback, seed := generateFeedback(10, 0, mutator)
	t.Logf("seed: %x\n", seed)

	for _, f := range expectedFeedback {
		expect := schema
		expect.FromCSVString(feedbackToCSV(f))

		mock.ExpectBegin()
		mock.
			ExpectQuery("^INSERT INTO [`\"']feedbacks[`\"'] .*$").
			WithArgs(f.Course, f.Feedback, 0, 0).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(f.ID))
		mock.ExpectCommit()

		actual, err := db.AddFeedback(f)
		assert.NoError(t, err, "there was an error adding the feedback to the database")
		assert.Equal(t, f, actual, "the expected and the actual feedback should be the same")
	}
}

// TestAddFeedbackInvalidCourse is a unit test that tests the
// return value and the stored data in the database using
// randomly created invalid feedback (by faker).
// AddFeedback is expected to reject invalid feedback, meaning
// feedback with missing Course.
func TestAddFeedbackInvalidCourse(t *testing.T) {
	db, closer, _, _ := createMockDatabase(t)
	defer closer()

	mutator := func(f *Feedback, fkr faker.Faker) {
		f.Course = ""
	}
	expectedFeedback, seed := generateFeedback(10, 0, mutator)
	t.Logf("seed: %x\n", seed)

	for _, f := range expectedFeedback {
		actual, err := db.AddFeedback(f)
		assert.Error(t, err, "invalid feedback should return an error")
		assert.ErrorIs(t, err, ErrInvalidFeedback, "add invalid feedback should return invalid feedback error")
		assert.Nil(t, actual, "no feedback should be returned when there is an error")
	}
}

// TestAddFeedbackInvalidFeedback is a unit test that tests the
// return value and the stored data in the database using
// randomly created invalid feedback (by faker).
// AddFeedback is expected to reject invalid feedback, meaning
// feedback with missing Feedback.
func TestAddFeedbackInvalidFeedback(t *testing.T) {
	db, closer, _, _ := createMockDatabase(t)
	defer closer()

	mutator := func(f *Feedback, fkr faker.Faker) {
		f.Feedback = ""
	}
	expectedFeedback, seed := generateFeedback(10, 0, mutator)
	t.Logf("seed: %x\n", seed)

	for _, f := range expectedFeedback {
		actual, err := db.AddFeedback(f)
		assert.Error(t, err, "invalid feedback should return an error")
		assert.ErrorIs(t, err, ErrInvalidFeedback, "add invalid feedback should return invalid feedback error")
		assert.Nil(t, actual, "no feedback should be returned when there is an error")
	}
}

// TestAddFeedbackError is a unit test that tests the return
// value and the stored data in the database using randomly
// created feedback (by faker).
// AddFeedback is expected to set the correct ID, and reset
// the Up- and Downvotes, add the sanitized data to the database
// and return this structure.
func TestAddFeedbackError(t *testing.T) {
	db, closer, mock, _ := createMockDatabase(t)
	defer closer()

	expectedFeedback, seed := generateFeedback(10, 0, nil)
	t.Logf("seed: %x\n", seed)

	for _, f := range expectedFeedback {
		mock.ExpectBegin()
		mock.
			ExpectQuery("^INSERT INTO [`\"']feedbacks[`\"'] .*$").
			WithArgs(f.Course, f.Feedback, 0, 0).
			WillReturnError(gorm.ErrDuplicatedKey)
		mock.ExpectRollback()

		actual, err := db.AddFeedback(f)
		assert.Error(t, err, "duplicate key should return an error")
		assert.ErrorIs(t, err, gorm.ErrDuplicatedKey, "add feedback with duplicate key should return duplicate key error")
		assert.Nil(t, actual, "no feedback should be returned when there is an error")
	}
}

// TestUpdateFeedback is a unit test that tests the return
// value and the stored data in the database using randomly
// created feedback.
// UpdateFeedback is expected to set the database values to the
// updated feedback passed, matching is based on ID, and return
// the updated data.
func TestUpdateFeedback(t *testing.T) {
	db, closer, mock, schema := createMockDatabase(t)
	defer closer()

	expectedFeedback, seed := generateFeedback(10, 0, nil)
	t.Logf("seed: %x\n", seed)

	for _, f := range expectedFeedback {
		expect := schema
		expect.FromCSVString(feedbackToCSV(f))

		mock.ExpectBegin()
		mock.
			ExpectExec("^UPDATE [`\"']feedbacks[`\"'] SET .* WHERE .*$").
			WithArgs(f.Course, f.Feedback, f.Upvotes, f.Downvotes, f.ID).
			WillReturnResult(sqlmock.NewResult(int64(f.ID), 1))
		mock.ExpectCommit()

		actual, err := db.UpdateFeedback(f)
		assert.NoError(t, err, "update valid feedback should not return an error")
		assert.Equal(t, f, actual, "returned feedback should be the same as the updated one")
	}
}

// TestUpdateFeedbackInvalidCourse is a unit test that tests
// the return value and the stored data in the database using
// randomly created feedback.
// UpdateFeedback is expected reject invalid feedback updates.
func TestUpdateFeedbackInvalidCourse(t *testing.T) {
	db, closer, _, _ := createMockDatabase(t)
	defer closer()

	mutator := func(f *Feedback, fkr faker.Faker) {
		f.Course = ""
	}
	expectedFeedback, seed := generateFeedback(10, 0, mutator)
	t.Logf("seed: %x\n", seed)

	for _, f := range expectedFeedback {
		actual, err := db.UpdateFeedback(f)
		assert.Error(t, err, "update invalid feedback should return invalid feedback error")
		assert.ErrorIs(t, err, ErrInvalidFeedback, "update invalid feedback should return invalid feedback error")
		assert.Nil(t, actual, "returned feedback should be nil on invalid update")
	}
}

// TestUpdateFeedbackInvalidFeedback is a unit test that tests
// the return value and the stored data in the database using
// randomly created feedback.
// UpdateFeedback is expected reject invalid feedback updates.
func TestUpdateFeedbackInvalidFeedback(t *testing.T) {
	db, closer, _, _ := createMockDatabase(t)
	defer closer()

	mutator := func(f *Feedback, fkr faker.Faker) {
		f.Feedback = ""
	}
	expectedFeedback, seed := generateFeedback(10, 0, mutator)
	t.Logf("seed: %x\n", seed)

	for _, f := range expectedFeedback {
		actual, err := db.UpdateFeedback(f)
		assert.Error(t, err, "update invalid feedback should return invalid feedback error")
		assert.ErrorIs(t, err, ErrInvalidFeedback, "update invalid feedback should return invalid feedback error")
		assert.Nil(t, actual, "returned feedback should be nil on invalid update")
	}
}

// TestUpdateFeedbackUnknown is a unit test that tests the return
// value and the stored data in the database using randomly
// created feedback.
// UpdateFeedback is expected to reject updates for unknown
// feedback.
func TestUpdateFeedbackUnknown(t *testing.T) {
	db, closer, mock, _ := createMockDatabase(t)
	defer closer()

	expectedFeedback, seed := generateFeedback(10, 0, nil)
	t.Logf("seed: %x\n", seed)

	for _, f := range expectedFeedback {
		mock.ExpectBegin()
		mock.
			ExpectExec("^UPDATE [`\"']feedbacks[`\"'] SET .* WHERE .*$").
			WithArgs(f.Course, f.Feedback, f.Upvotes, f.Downvotes, f.ID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		actual, err := db.UpdateFeedback(f)
		assert.Error(t, err, "non existent feedback should return no such feedback error")
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound, "non existent feedback should return no such feedback error")
		assert.Nil(t, actual, "non existent feedback should not return any feedback")
	}
}

// TestDeleteFeedback is a unit test that tests the deletion
// or voiding of data in the database. A void feedback is a feedback
// with a 0 ID. A deleted feedback is completely removed from the
// database. Feedback is generated randomly,
// DeleteFeedback is expected to delete the provided feedback if
// it exists, or return an error if it doesn't. Matching is done
// using ALL attributes, so that mismatching versions of a feedback
// are considered.
func TestDeleteFeedback(t *testing.T) {
	db, closer, mock, _ := createMockDatabase(t)
	defer closer()

	expectedFeedback, seed := generateFeedback(10, 0, nil)
	t.Logf("seed: %x\n", seed)

	for _, f := range expectedFeedback {
		mock.ExpectBegin()
		mock.
			ExpectExec("^DELETE FROM [`\"']feedbacks[`\"] .*$").
			WithArgs(f.Course, f.Feedback, f.Upvotes, f.Downvotes, f.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := db.DeleteFeedback(f)
		assert.NoError(t, err, "deleting an existing feedback should not be a problem")
	}
}

// TestDeleteFeedbackUnknown is a unit test that tests the deletion
// or voiding of data in the database. A void feedback is a feedback
// with a 0 ID. A deleted feedback is completely removed from the
// database. Feedback is generated randomly,
// DeleteFeedback is expected to delete the provided feedback if
// it exists, or return an error if it doesn't. Matching is done
// using ALL attributes, so that mismatching versions of a feedback
// are considered.
func TestDeleteFeedbackUnknown(t *testing.T) {
	db, closer, mock, _ := createMockDatabase(t)
	defer closer()

	expectedFeedback, seed := generateFeedback(10, 0, nil)
	t.Logf("seed: %x\n", seed)

	for _, f := range expectedFeedback {
		mock.ExpectBegin()
		mock.
			ExpectExec("^DELETE FROM [`\"']feedbacks[`\"] .*$").
			WithArgs(f.Course, f.Feedback, f.Upvotes, f.Downvotes, f.ID).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		err := db.DeleteFeedback(f)
		assert.Error(t, err, "deleting unknown feedback should return an error")
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound, "deleting unknown feedback should return not found error")
	}
}

// TestIncrementFeedbackUpvotes tests the correct incrementing
// of feedback upvotes.
// A feedback should be atomically updated to include the
// upvote. Upvotes may not exceed 2000 or go below 0.
func TestIncrementFeedbackUpvotes(t *testing.T) {
	db, closer, mock, schema := createMockDatabase(t)
	defer closer()

	expectedFeedback, seed := generateFeedback(10, 0, func(f *Feedback, fkr faker.Faker) {
		f.Upvotes = fkr.UIntBetween(0, 1999)
	})

	t.Logf("seed: %x\n", seed)

	for _, f := range expectedFeedback {
		expectSQL := schema
		expectSQL.FromCSVString(feedbackToCSV(f))
		f.Upvotes += 1

		mock.ExpectBegin()
		mock.
			ExpectQuery("^SELECT .+ FROM [`\"']feedbacks[`\"'] WHERE [`\"']feedbacks[`\"']\\.[`\"']id[`\"']\\W*=.*$").
			WithArgs(f.ID, 1).
			WillReturnRows(expectSQL)
		mock.
			ExpectExec("^UPDATE [`\"']feedbacks[`\"'] SET .* WHERE .*$").
			WithArgs(f.Course, f.Feedback, f.Upvotes, f.Downvotes, f.ID).
			WillReturnResult(sqlmock.NewResult(int64(f.ID), 1))
		mock.ExpectCommit()

		actual, err := db.IncrementFeedbackUpvotes(f.ID)
		assert.NoError(t, err, "update valid feedback should not return an error")
		assert.Equal(t, f, actual, "returned feedback should be the same as the updated one")
	}
}

// TestDecrementFeedbackUpvotes tests the correct decrementing
// of feedback upvotes.
// A feedback should be atomically updated to include the
// upvote. Upvotes may not exceed 2000 or go below 0.
func TestDecrementFeedbackUpvotes(t *testing.T) {
	db, closer, mock, schema := createMockDatabase(t)
	defer closer()

	expectedFeedback, seed := generateFeedback(10, 0, func(f *Feedback, fkr faker.Faker) {
		f.Upvotes = fkr.UIntBetween(1, 2000)
	})
	t.Logf("seed: %x\n", seed)

	for _, f := range expectedFeedback {
		expectSQL := schema
		expectSQL.FromCSVString(feedbackToCSV(f))
		f.Upvotes -= 1

		mock.ExpectBegin()
		mock.
			ExpectQuery("^SELECT .+ FROM [`\"']feedbacks[`\"'] WHERE [`\"']feedbacks[`\"']\\.[`\"']id[`\"']\\W*=.*$").
			WithArgs(f.ID, 1).
			WillReturnRows(expectSQL)
		mock.
			ExpectExec("^UPDATE [`\"']feedbacks[`\"'] SET .* WHERE .*$").
			WithArgs(f.Course, f.Feedback, f.Upvotes, f.Downvotes, f.ID).
			WillReturnResult(sqlmock.NewResult(int64(f.ID), 1))
		mock.ExpectCommit()

		actual, err := db.DecrementFeedbackUpvotes(f.ID)
		assert.NoError(t, err, "update valid feedback should not return an error")
		assert.Equal(t, f, actual, "returned feedback should be the same as the updated one")
	}
}

// TestIncrementFeedbackDownvotes tests the correct incrementing
// of feedback upvotes.
// A feedback should be atomically updated to include the
// upvote. Upvotes may not exceed 2000 or go below 0.
func TestIncrementFeedbackDownvotes(t *testing.T) {
	db, closer, mock, schema := createMockDatabase(t)
	defer closer()

	expectedFeedback, seed := generateFeedback(10, 0, func(f *Feedback, fkr faker.Faker) {
		f.Downvotes = fkr.UIntBetween(0, 1999)
	})

	t.Logf("seed: %x\n", seed)

	for _, f := range expectedFeedback {
		expectSQL := schema
		expectSQL.FromCSVString(feedbackToCSV(f))
		f.Downvotes += 1

		mock.ExpectBegin()
		mock.
			ExpectQuery("^SELECT .+ FROM [`\"']feedbacks[`\"'] WHERE [`\"']feedbacks[`\"']\\.[`\"']id[`\"']\\W*=.*$").
			WithArgs(f.ID, 1).
			WillReturnRows(expectSQL)
		mock.
			ExpectExec("^UPDATE [`\"']feedbacks[`\"'] SET .* WHERE .*$").
			WithArgs(f.Course, f.Feedback, f.Upvotes, f.Downvotes, f.ID).
			WillReturnResult(sqlmock.NewResult(int64(f.ID), 1))
		mock.ExpectCommit()

		actual, err := db.IncrementFeedbackDownvotes(f.ID)
		assert.NoError(t, err, "update valid feedback should not return an error")
		assert.Equal(t, f, actual, "returned feedback should be the same as the updated one")
	}
}

// TestDecrementFeedbackDownvotes tests the correct decrementing
// of feedback Downvotes.
// A feedback should be atomically updated to include the
// upvote. Downvotes may not exceed 2000 or go below 0.
func TestDecrementFeedbackDownvotes(t *testing.T) {
	db, closer, mock, schema := createMockDatabase(t)
	defer closer()

	expectedFeedback, seed := generateFeedback(10, 0, func(f *Feedback, fkr faker.Faker) {
		f.Downvotes = fkr.UIntBetween(1, 2000)
	})
	t.Logf("seed: %x\n", seed)

	for _, f := range expectedFeedback {
		expectSQL := schema
		expectSQL.FromCSVString(feedbackToCSV(f))
		f.Downvotes -= 1

		mock.ExpectBegin()
		mock.
			ExpectQuery("^SELECT .+ FROM [`\"']feedbacks[`\"'] WHERE [`\"']feedbacks[`\"']\\.[`\"']id[`\"']\\W*=.*$").
			WithArgs(f.ID, 1).
			WillReturnRows(expectSQL)
		mock.
			ExpectExec("^UPDATE [`\"']feedbacks[`\"'] SET .* WHERE .*$").
			WithArgs(f.Course, f.Feedback, f.Upvotes, f.Downvotes, f.ID).
			WillReturnResult(sqlmock.NewResult(int64(f.ID), 1))
		mock.ExpectCommit()

		actual, err := db.DecrementFeedbackDownvotes(f.ID)
		assert.NoError(t, err, "update valid feedback should not return an error")
		assert.Equal(t, f, actual, "returned feedback should be the same as the updated one")
	}
}
