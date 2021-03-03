package postgres

import (
	"errors"
	"time"

	"github.com/form3tech-oss/jwt-go"
)

const (
	Issuer     = "authorize"
	SecrateKey = "secerteKey"
)

type jwtClaim struct {
	UserId string
	jwt.StandardClaims
}

func GenrateToken(userId string, password string) (string, error) {
	jt := jwtClaim{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			Issuer:    Issuer,
		},
	}
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, jt)
	token, err := tk.SignedString([]byte(SecrateKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func TokenValidation(token string) (*jwtClaim, error) {

	jtoken, err := jwt.ParseWithClaims(token, &jwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecrateKey), nil
		})
	if err != nil {
		log.Errorf("Failed to parse token %s", err.Error())
		return nil, err
	}
	claim, ok := jtoken.Claims.(*jwtClaim)
	if !ok {
		return nil, errors.New("Failed to pasrse claim")
	}

	if claim.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("Token is expired")
	}
	return claim, nil
}
