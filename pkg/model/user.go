package model

import (
	"github.com/supergiant/supergiant/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserRoleAdmin = "admin"
	UserRoleUser  = "user"
)

type UserList struct {
	Pagination
	Items []*User `json:"items"`
}

type User struct {
	BaseModel

	Username string `json:"username" validate:"nonzero,max=24,regexp=^[A-Za-z0-9_-]+$" gorm:"not null;unique_index"`
	Password string `json:"password,omitempty" validate:"nonzero,min=8,max=32" gorm:"-" sg:"private"`
	Role     string `json:"role" validate:"nonzero" gorm:"not null" sg:"default=user"`

	EncryptedPassword []byte `json:"-" gorm:"not null"`

	APIToken string `json:"api_token" gorm:"not null;index" sg:"readonly"`
}

func (m *User) BeforeCreate() error {
	m.GenerateAPIToken()
	return nil
}

func (m *User) BeforeSave() error {
	if m.Password == "" {
		return nil
	}
	return m.encryptPassword()
}

func (m *User) GenerateAPIToken() {
	m.APIToken = util.RandomString(32)
}

////////////////////////////////////////////////////////////////////////////////
// Private                                                                    //
////////////////////////////////////////////////////////////////////////////////

func (m *User) encryptPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.EncryptedPassword = hashedPassword
	m.Password = "" // just for extra-good measure
	return nil
}
