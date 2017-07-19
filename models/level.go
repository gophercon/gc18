package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

type Level struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
}

// String is not required by pop and may be deleted
func (l Level) String() string {
	jl, _ := json.Marshal(l)
	return string(jl)
}

// Levels is not required by pop and may be deleted
type Levels []Level

// String is not required by pop and may be deleted
func (l Levels) String() string {
	jl, _ := json.Marshal(l)
	return string(jl)
}

// Validate gets run every time you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (l *Level) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: l.Name, Name: "Name"},
		&validators.StringIsPresent{Field: l.Description, Name: "Description"},
	), nil
}

// ValidateSave gets run every time you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (l *Level) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (l *Level) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
