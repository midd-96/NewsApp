package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/upper/db/v4"
)

var (
	ErrNoMoreRows     = errors.New("No records found")
	ErrDuplicateEmail = errors.New("Email already exists")
	ErrUserNotActive  = errors.New("Your account is inactive")
	ErrInvalidLogin   = errors.New("Invalid Username/password")
)

type Models struct {
	Users UserModel
}

func New(db db.Session) Models {
	return Models{
		Users: UserModel{db: db},
	}
}

func convertUpperIDtoInt(id db.ID) int {
	idType := fmt.Sprintf("%T", id)
	if idType == "int64" {
		return int(id.(int64))
	}
	return id.(int)
}

func errHasDuplicate(err error, key string) bool {
	str := fmt.Sprintf(`ERROR: duplicate key value violate unique constraint "%s"`, key)
	return strings.Contains(err.Error(), str)
}
