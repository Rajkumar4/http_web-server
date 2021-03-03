package handler

import (
	"database/sql"
	"fmt"

	pg "github.com/httpWebServer/m/postgres"
	logger "github.com/ipfs/go-log/v2"
)

var log = logger.Logger("handler")

func SignUp(email string, name string, password string, db *sql.DB) (error, int64) {
	usr := &pg.User{
		Email:    email,
		Name:     name,
		Password: password,
	}
	err, id := pg.Insert(usr, db)
	if err != nil {
		log.Errorf("Failed to register user %s", err.Error())
		return err, id
	}
	return nil, id
}

func Login(userId string, password string, db *sql.DB) (string, error) {
	cdr := &pg.Crdentials{
		UserId:   userId,
		Password: password,
	}
	 err := pg.Read(db, cdr)
	if err != nil {
		log.Errorf("User doesn't exist %s", err.Error())
		return "", err
	}
	token, err := pg.GenrateToken(userId, password)
	if err != nil {
		log.Errorf("Failed to generate toke %s", err.Error())
		return "", err
	}
	usr := &pg.User{
		Email: userId,
		Token: token,
	}
	err = pg.Update(db, usr)
	if err != nil {
		log.Errorf("Failed to get update %s", err.Error())
		return "", err
	}
	return token, nil
}

func TokenValidation(token string, db *sql.DB) (string, error) {
	claim, err := pg.TokenValidation(token)
	if err != nil {
		log.Errorf("token is not valid %s", err.Error())
		return "", err
	}
	message := fmt.Sprintf("Welcome %s ", claim.UserId)
	return message, nil
}
