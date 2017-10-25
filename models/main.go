package models

import (
	"gopkg.in/mgo.v2"
)

var db *mgo.Database

type Dao struct {
	session *mgo.Session
}

// set mongodb
func NewDao() (*Dao, error) {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	db = session.DB("lilac-go")
	return &Dao{session}, nil
}

func (d *Dao) Close() {
	d.session.Close()
}
