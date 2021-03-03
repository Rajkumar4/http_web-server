package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	logger "github.com/ipfs/go-log/v2"
	_ "github.com/lib/pq"
)

var log = logger.Logger("postgres")

type Crdentials struct {
	UserId   string
	Password string
}

type User struct {
	Id       string
	Email    string
	Name     string
	Password string
	Token    string
}

const (
	host      = "localhost"
	post      = "8050"
	user      = "postgres"
	password  = "pgsql"
	dbname    = "testdb"
	tablename = "webserver"
)

func PostgresConfig() *sql.DB {
	pgInfo := fmt.Sprintf("user=%s dbname=%s", user, dbname)
	db, err := sql.Open("postgres", pgInfo)
	if err != nil {
		log.Errorf("Failed to connect postgres %s", err.Error())
		return nil
	}
	err = db.Ping()
	if err != nil {
		log.Errorf("ping failed %s", err.Error())
		return nil
	}
	err = tableCreate(db)
	if err != nil {
		log.Errorf("failed to create table %s", err.Error())
		return nil
	}
	return db
}

func tableCreate(db *sql.DB) error {
	query := fmt.Sprintf(`
	create table if not exists %s 
	(id serial,
	email varchar(50) PRIMARY KEY,
	name varchar(50),
	password varchar(50),
	token varchar(500)
	)`, tablename)
	_, err := db.Query(query)
	if err != nil {
		log.Errorf("Failed to preapre sql query %s", err.Error())
		return err
	}
	return nil
}

func Insert(usr *User, db *sql.DB) (error, int64) {
	query := fmt.Sprintf(`
	insert into %s 
	(email,
     name,
	 password)VALUES($1,$2,$3) returning id
	`, tablename)
	var id int64
	err := db.QueryRow(query, usr.Email, usr.Name, usr.Password).Scan(&id)
	if err != nil {
		log.Errorf("Failed to insert data %s", err.Error())
		return err, 0
	}
	return nil, id
}

func Read(db *sql.DB, cdr *Crdentials) error {
	query := fmt.Sprintf(`SELECT id, email,name, password from %s where email=$1`, tablename)
	rows, err := db.Query(query, cdr.UserId)
	if err != nil {
		log.Errorf("failed to read data %s", err.Error())
		return err
	}
	var usr User
	if rows.Next() {
		err = rows.Scan(&usr.Id, &usr.Email, &usr.Name, &usr.Password)
		if err != nil {
			log.Errorf("Failed to get values %s", err.Error())
			return err
		}
	}
	if cdr.UserId != usr.Email || cdr.Password != usr.Password {
		return errors.New("Invaild userid or password")
	}
	return nil
}

func Update(db *sql.DB, usr *User) error {
	query := fmt.Sprintf(`
	update %s SET token =$1 where email=$2
	`, tablename)
	result, err := db.Exec(query, usr.Token, usr.Email)
	if err != nil {
		log.Errorf("Failed to update table %s", err.Error())
		return nil
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Errorf("Failed to update record %s", err.Error())
		return err
	}
	return nil
}
