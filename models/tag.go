package models

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/Sirupsen/logrus"
)

//Tag struct contains post tag info
type Tag struct {
	Name string `form:"name" json:"name"`
	//calculated fields
	PostCount int64 `json:"post_count"`
}

//Insert saves tag info into db
func (tag *Tag) Insert() error {
	collection := db.C("tag")
	err := collection.Insert(tag)
	return err
}

//Delete removes tag from db, according to postgresql constraints tag associations are removed either
func (tag *Tag) Delete() error {
	collection := db.C("tag")
	err := collection.Remove(bson.M{"name": tag.Name})
	return err
}

//GetTag returns tag by its name
func GetTag(name interface{}) (*Tag, error) {
	tag := &Tag{}
	collection := db.C("tag")
	err := collection.Find(bson.M{"name": name}).One(&tag)
	return tag, err
}

//GetTags returns a slice of tags, ordered by name
func GetTags() ([]Tag, error) {
	var list []Tag
	collection := db.C("tag")
	err := collection.Find(nil).Sort("-_id").All(&list)
	logrus.Errorf("tag list : %s", list)
	return list, err
}

//UpdateTags inserts new (non existent) post tags and updates associations
func (post *Post) UpdateTags() error {
	collection := db.C("tag")
	for _, name := range post.Tags {
		post_count, _ := db.C("post").Find(bson.M{"tags": bson.M{"$in": []string{name}}}).Count()
		collection.Upsert(bson.M{"name": name}, bson.M{"name": name, "post_count": post_count})
	}
	return nil
}

//GetNotEmptyTags returns a slice of tags that have at least one associated blog post
func GetNotEmptyTags() ([]Tag, error) {
	var list []Tag
	collection := db.C("tag")
	err := collection.Find(nil).Sort("-_id").All(&list)
	return list, err
}
