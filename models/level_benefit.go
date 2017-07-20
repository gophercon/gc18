package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

type LevelBenefit struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	LevelID     string    `json:"level_id" db:"level_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
}

// String is not required by pop and may be deleted
func (l LevelBenefit) String() string {
	jl, _ := json.Marshal(l)
	return string(jl)
}

// LevelBenefits is not required by pop and may be deleted
type LevelBenefits []LevelBenefit

// String is not required by pop and may be deleted
func (l LevelBenefits) String() string {
	jl, _ := json.Marshal(l)
	return string(jl)
}

// Validate gets run every time you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (l *LevelBenefit) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: l.LevelID, Name: "LevelID"},
		&validators.StringIsPresent{Field: l.Name, Name: "Name"},
		&validators.StringIsPresent{Field: l.Description, Name: "Description"},
	), nil
}

// ValidateSave gets run every time you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (l *LevelBenefit) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (l *LevelBenefit) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
