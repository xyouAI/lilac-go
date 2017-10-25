package models

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/Sirupsen/logrus"
)

//Tag struct contains post tag info
type Category struct {
	Name string `form:"name" json:"name"`
	//calculated fields
	PostCount int64 `json:"post_count"`
}

//Insert saves tag info into db
func (category *Category) Insert() error {
	collection := db.C("category")
	err := collection.Insert(category)
	return err
}

//Delete removes tag from db, according to postgresql constraints tag associations are removed either
func (category *Category) Delete() error {
	collection := db.C("category")
	err := collection.Remove(bson.M{"name": category.Name})
	return err
}

//GetTag returns tag by its name
func GetCategory(name interface{}) (*Category, error) {
	category := &Category{}
	collection := db.C("category")
	err := collection.Find(bson.M{"name": name}).One(&category)
	return category, err
}

//GetTags returns a slice of tags, ordered by name
func GetCategorys() ([]Category, error) {
	var list []Category
	collection := db.C("category")
	err := collection.Find(nil).Sort("-_id").All(&list)
	logrus.Errorf("category list : %s", list)
	return list, err
}

//UpdateTags inserts new (non existent) post tags and updates associations
func (post *Post) UpdateCategorys() error {
	//collection := db.C("category")
	//post_count,_ := db.C("post").Find(bson.M{"category":bson.M{"$in":[]string{post.Category}}}).Count()
	//collection.Upsert(bson.M{"name":post.Category},bson.M{"name":post.Category,"post_count":post_count})
	return nil
}
