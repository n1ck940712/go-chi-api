package models

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        int32
	Username  string `gorm:"text;size:255;not null;unique;default:null"`
	Password  string `gorm:"text;size:100;not null;default:null"`
	Email     string `gorm:"text;null;default:null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Username == "" {
			return errors.New("required username")
		}
		if u.Password == "" {
			return errors.New("required password")
		}
		// if u.Email != "" {
		// 	if err := checkmail.ValidateFormat(u.Email); err != nil {
		// 		return errors.New("Invalid Email")
		// 	}
		// }
		return nil

	case "login":
		if u.Username == "" {
			return errors.New("required username")
		}
		if u.Password == "" {
			return errors.New("required password")
		}
		return nil

	default:
		if u.Username == "" {
			return errors.New("required username")
		}
		if u.Password == "" {
			return errors.New("required password")
		}
		// if u.Email != "" {
		// 	if err := checkmail.ValidateFormat(u.Email); err != nil {
		// 		return errors.New("Invalid Email")
		// 	}
		// }
		return nil
	}
}

func (u *User) GetAll(db *gorm.DB) (*[]User, error) {
	users := []User{}
	err := db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) Get(db *gorm.DB, id int32) (*User, error) {
	err := db.Debug().Model(User{}).Where("id = ?", id).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

func (u *User) Save(db *gorm.DB) (*User, error) {
	u.Password = HashPassword(u.Password)
	err := db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) Update(db *gorm.DB, id int32) (*User, error) {
	HashPassword(u.Password)
	db = db.Debug().Model(&User{}).Where("id = ?", id).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"nickname":  u.Username,
			"password":  u.Password,
			"email":     u.Email,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	err := db.Debug().Model(&User{}).Where("id = ?", id).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) Delete(db *gorm.DB, id int32) (int64, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", id).Take(&User{}).Delete(&User{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
