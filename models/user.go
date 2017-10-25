package models

import (
	//"fmt"
	"time"
	"golang.org/x/crypto/bcrypt"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus"
)

//User type contains user info
type User struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Email     string        `form:"email" json:"email"`
	Name      string        `form:"name" json:"name"`
	Password  string        `form:"password" json:"password" bson:",omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}

//HashPassword substitutes User.Password with its bcrypt hash
func (user *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return nil
}

//Insert saves user info into db
func (user *User) Insert() error {
	//userData := User{
	//	email:user.Email,
	//}
	collection := db.C("user")
	err := collection.Insert(user)
	//err := db.QueryRow("INSERT INTO users(email, name, password, timestamp) VALUES(lower($1),$2,$3,$4) RETURNING id", user.Email, user.Name, user.Password, time.Now()).Scan(&user.ID)
	return err
}

//Update updates user info in db
func (user *User) Update(id string) error {
	collection := db.C("user")
	user.Timestamp = time.Now()
	objectId := bson.ObjectIdHex(id)
	err := collection.UpdateId(objectId, user)
	//_, err := db.Exec("UPDATE users SET email=lower($2), name=$3, password=$4 WHERE id=$1", user.ID, user.Email, user.Name, user.Password)
	return err
}

//Delete removes user record from db. Can't remove the last user
func (user *User) Delete(id string) error {
	collection := db.C("user")
	objectId := bson.ObjectIdHex(id)
	err := collection.RemoveId(objectId)
	return err
}

//GetUser returns user by his id
func GetUser(id interface{}) (*User, error) {
	user := &User{}
	collection := db.C("user")
	//objectId := bson.ObjectIdHex(id)
	err := collection.FindId(id).One(&user)
	logrus.Error("getuser", id, user)
	//err := db.Get(user, "SELECT * FROM users WHERE id=$1", id)
	return user, err
}

//GetUsers returns a list of user ordered by id
func GetUsers() ([]User, error) {
	var list []User
	collection := db.C("user")
	err := collection.Find(nil).All(&list)
	//err := db.Select(&list, "SELECT * FROM users ORDER BY id")
	return list, err
}

//GetUserByEmail returns user record by his email, case insensitive
func GetUserByEmail(email interface{}) (*User, error) {
	user := &User{}
	collection := db.C("user")
	err := collection.Find(bson.M{"email": email}).One(&user)
	//err := db.Get(user, "SELECT * FROM users WHERE lower(email)=lower($1)", email)
	return user, err
}
