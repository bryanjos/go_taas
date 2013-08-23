package main

import (
	"errors"
	r "github.com/christopherhesse/rethinkgo"
)

type DB struct {
	session   *r.Session
	tableName string
}

func (db DB) find(key string, value string, i interface{}) error {
	rows := r.Table(db.tableName).Filter(r.Row.Attr(key).Eq(value)).Run(db.session)

	isSet := false
	for rows.Next() {
		if err := rows.Scan(&i); err != nil {
			return err
		}

		isSet = true
	}

	if isSet {
		return nil
	} else {
		return errors.New("Not Found")
	}

	return errors.New("Database Error")
}

func (db DB) create(i interface{}) {
	r.Table(db.tableName).Insert(i).Run(db.session).Exec()
}

func (db DB) update(i interface{}) {
	r.Table(db.tableName).Update(i).Run(db.session).Exec()
}

func (db DB) delete(key string, value string) {
	r.Table(db.tableName).Filter(r.Row.Attr(key).Eq(value)).Delete().Run(db.session).Exec()
}
