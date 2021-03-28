package model

import (
	"errors"
	"strings"
	"time"
)

// User represents an User
type User struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Nick      string    `json:"nick,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// Prepare call methods to validate and format the data of user
func (user *User) Prepare(step string) error {
	if err := user.validate(step); err != nil {
		return err
	}

	user.format()
	return nil
}

func (user *User) validate(step string) error {
	if user.Name == "" {
		return errors.New("o campo nome é obrigatório")
	}

	if user.Nick == "" {
		return errors.New("o campo nick é obrigatório")
	}

	if user.Email == "" {
		return errors.New("o campo email é obrigatório")
	}

	if step == "register" && user.Password == "" {
		return errors.New("o campo senha é obrigatório")
	}

	return nil
}

func (user *User) format() {
	user.Name = strings.TrimSpace(user.Name)
	user.Nick = strings.TrimSpace(user.Nick)
	user.Email = strings.TrimSpace(user.Email)
}
