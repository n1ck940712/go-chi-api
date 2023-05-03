package auth

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"go-chi-api/internal/database"
	"go-chi-api/internal/models"

	jwt "github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

const (
	TokenExpiryDuration = time.Hour * 24 * 1 // 1 day
)

func CreateToken(user *models.User) (string, error) {
	// Generate a 32-byte random string
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// Encode the random string in base64
	token := sha512.Sum512(b)
	tokenString := hex.EncodeToString(token[:])

	// Save the session token in the database
	newSession := models.SessionTable{
		UserID:  user.ID,
		Token:   tokenString,
		Expires: time.Now().Add(TokenExpiryDuration),
	}

	oldSession, err := models.GetSessionByUserID(user.ID)

	var result *gorm.DB

	if err != nil {
		result = database.DB.Create(&newSession)
	} else {
		result = database.DB.Model(&models.SessionTable{}).Where("id = ?", oldSession.ID).Updates(newSession)
	}
	if result.Error != nil {
		return "", result.Error
	}
	if result.RowsAffected == 0 {
		return "", errors.New("failed to create session token")
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) (*models.User, error) {
	tokenString, err := ExtractToken(r)
	if err != nil {
		return &models.User{}, err
	}
	session, err := models.GetSessionByToken(tokenString)
	fmt.Println(session)
	if err != nil {
		return &models.User{}, err
	}
	user2 := models.User{}
	user, err := user2.Get(database.DB, session.UserID)
	if err != nil {
		return &models.User{}, err
	}
	return user, nil
}

func TokenValidIsAdmin(r *http.Request) error {
	tokenString, err := ExtractToken(r)
	if err != nil {
		return err
	}
	session, err := models.GetSessionByToken(tokenString)
	if err != nil {
		return err
	}
	user := models.User{}
	fetchedUser, err := user.Get(database.DB, session.UserID)
	if err != nil {
		return err
	}

	if fetchedUser.Role != "admin" {
		fmt.Printf("not admin: %s", session.User.Role)
		return errors.New("unauthorized")
	}
	return nil
}

func ExtractToken(r *http.Request) (string, error) {
	bearToken := r.URL.Query().Get("token")
	if bearToken != "" {
		return bearToken, nil
	}
	bearToken = r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1], nil
	}
	return "", nil
}
