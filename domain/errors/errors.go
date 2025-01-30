package errors

import "errors"

// Common errors
var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("resource not found")
)

// Word related errors
var (
	ErrEmptyWordText       = errors.New("word text cannot be empty")
	ErrEmptyExample        = errors.New("example cannot be empty")
	ErrEmptyTag            = errors.New("tag cannot be empty")
	ErrWordNotFound        = errors.New("word not found")
	ErrWordAlreadyExists   = errors.New("word already exists")
	ErrInvalidDifficulty   = errors.New("difficulty must be between 1 and 5")
	ErrInvalidMasteryLevel = errors.New("mastery level must be between 0 and 5")
	ErrDuplicateTag        = errors.New("tag already exists")
)

// Translation related errors
var (
	ErrEmptyPrimaryTranslation     = errors.New("primary translation cannot be empty")
	ErrInvalidSecondaryTranslation = errors.New("invalid secondary translation")
	ErrEmptyTranslation            = errors.New("translation cannot be empty")
)

// Vocabulary related errors
var (
	ErrEmptyVocabularyName = errors.New("vocabulary name cannot be empty")
	ErrVocabularyNotFound  = errors.New("vocabulary not found")
	ErrDuplicateWord       = errors.New("word already exists in vocabulary")
)

// Repository related errors
var (
	ErrFailedToSave   = errors.New("failed to save")
	ErrFailedToUpdate = errors.New("failed to update")
	ErrFailedToDelete = errors.New("failed to delete")
	ErrFailedToQuery  = errors.New("failed to query")
)

// Learning progress related errors
var (
	ErrProgressNotFound = errors.New("learning progress not found")
)

// User related errors
var (
	ErrInvalidUserID = errors.New("invalid user ID")
)
