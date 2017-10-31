package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"github.com/Sirupsen/logrus"
)

type Page struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
	Slug           string        `form:"slug" json:"slug"`
	Name           string        `form:"name" json:"name"`
	Description    string        `form:"page-editormd-html-code" json:"page-editormd-html-code"`
	DescriptionDoc string        `form:"page-editormd-markdown-doc" json:"page-editormd-markdown-doc"`
	Published      bool          `form:"published" json:"published"`
	CreatedAt      time.Time     `form:"_" json:"created_at" db:"created_at"`
	UpdatedAt      time.Time     `form:"_" json:"updated_at" db:"updated_at"`
}

type Pages struct {
	page []Page
}

//Insert saves Page struct into db
func (page *Page) Insert() error {
	collection := db.C("page")
	err := collection.Insert(page)
	return err
}

//Update saves Page changes into db
func (page *Page) Update(id string) error {
	collection := db.C("page")
	objectId := bson.ObjectIdHex(id)
	page.UpdatedAt = time.Now()
	err := collection.Update(bson.M{"_id": objectId}, bson.M{"$set": page})
	return err
}

//Delete removes page from db
func (page *Page) Delete(id string) error {
	collection := db.C("page")
	objectId := bson.ObjectIdHex(id)
	err := collection.RemoveId(objectId)
	return err
}

//GetPage fetches page from db by its id
func GetPage(id string) (*Page, error) {
	page := &Page{}
	collection := db.C("page")
	objectId := bson.ObjectIdHex(id)
	err := collection.Find(bson.M{"_id": objectId}).One(&page)
	return page, err
}

func GetPageBySlug(slug string) (*Page, error) {
	page := &Page{}
	collection := db.C("page")
	err := collection.Find(bson.M{"slug": slug}).One(&page)
	logrus.Error(slug,page)
	return page, err
}

//GetPages returns a slice of all pages
func GetPages() ([]Page, error) {
	var list []Page
	collection := db.C("page")
	err := collection.Find(nil).All(&list)
	return list, err
}

//GetPublishedPages returns a slice of pages with .Published=true
func GetPublishedPages() ([]Page, error) {
	var list []Page
	collection := db.C("page")
	err := collection.Find(nil).One(&list)
	//err := db.Select(&list, "SELECT * FROM pages WHERE published=$1 ORDER BY id", true)
	return list, err
}
