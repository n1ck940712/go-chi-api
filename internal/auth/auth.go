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
	session := models.SesssionTable{
		UserID:  user.ID,
		Token:   tokenString,
		Expires: time.Now().Add(TokenExpiryDuration),
	}
	result := database.DB.Create(&session)
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

func TokenValid(r *http.Request) error {
	tokenString, err := ExtractToken(r)
	if err != nil {
		return err
	}
	session, err := models.GetSessionByToken(tokenString)
	if err != nil {
		return err
	}
	return nil
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
