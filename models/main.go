package models

import (
	"gopkg.in/mgo.v2"
	"github.com/phpfor/lilac-go/system"
	"fmt"
)

var db *mgo.Database

type Dao struct {
	session *mgo.Session
}

// set mongodb
func NewDao() (*Dao, error) {
	c := system.GetConfig().Database
	url := fmt.Sprintf("%s:%s",c.Host,c.Port)
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	db = session.DB(c.Name)
	return &Dao{session}, nil
}

func (d *Dao) Close() {
	d.session.Close()
}
