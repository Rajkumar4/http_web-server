package postgres

import (
	"testing"
)

func TestPostgres(t *testing.T) {
	db := PostgresConfig()
	if db == nil {
		t.Fatal()
	}
	cdr := &Crdentials{
		UserId:   "bootstrap@gmail.com",
		Password: "abcd@1234",
	}
	err := Read(db, cdr)
	if err != nil {
		t.Fatal()
	}
}
