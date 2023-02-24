package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func CreateToken(userId int32) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	fmt.Printf("claims: %v", claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
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
	fmt.Printf("request tokenString: %s", tokenString)
	if err != nil {
		return err
	}
	token, err := ParseToken(tokenString)
	if err != nil {
		return err
	}
	if !token.Valid {
		return err
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

func ExtractTokenID(r *http.Request) (int32, error) {
	token, err := ExtractToken(r)
	if err != nil {
		return 0, err
	}
	claims, err := ExtractTokenMetadata(token)
	if err != nil {
		return 0, err
	}
	return claims["user_id"].(int32), nil
}

func ExtractTokenMetadata(tokenString string) (jwt.MapClaims, error) {
	token, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
