package models

import (
	"html/template"
	//"github.com/jmoiron/sqlx"
	//"github.com/microcosm-cc/bluemonday"
	//"github.com/russross/blackfriday"
	//"gopkg.in/guregu/null.v3"
	//"fmt"
	"gopkg.in/mgo.v2/bson"
	//"github.com/revel/modules/db/app"
	"time"
	"github.com/Sirupsen/logrus"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

//Post type contains blog post info
type Post struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Title       string        `form:"title" json:"title"`
	Slug        string        `form:"slug" json:"slug"`
	Category    string        `form:"category" json:"category"`
	Description string        `form:"description" json:"description"`
	Published   bool          `form:"published" json:"published"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at,omitempty"`
	Author      User          `json:"author"`
	Tags        []string      `form:"tags" json:"tags" bson:",omitempty"` //can't make gin Bind form field to []Tag, so use []string instead
	//CommentCount int64    `form:"-" json:"comment_count" db:"comment_count"`
}

//type TagList struct {
//	Tags []string `form:"tags" bson:"tags"`
//}

//Insert saves Post as well as associated tags (creating them if needed) into db. Obsolete associations are removed
func (post *Post) Insert() error {
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	collection := db.C("post")
	err := collection.Insert(post)
	if err == nil {
		post.UpdateTags()
	}
	return err
}

//Update saves Post and associated tags changes into db
func (post *Post) Update(id string) error {
	collection := db.C("post")
	logrus.Println("post edit", post)
	objectId := bson.ObjectIdHex(id)
	post.UpdatedAt = time.Now()
	logrus.Println("objectId : ", objectId)
	err := collection.Update(bson.M{"_id": objectId}, bson.M{"$set": post})
	if err == nil {
		post.UpdateTags()
	}

	return err
}

//Delete removes Post from db. Existing postgresql contstraints remove tag associations on post delete.
func (post *Post) Delete(id string) error {
	collection := db.C("post")
	objectId := bson.ObjectIdHex(id)
	err := collection.RemoveId(objectId)
	return err
}

//Excerpt returns post excerpt by removing html tags first and truncating to 300 symbols
func (post *Post) Excerpt() template.HTML {
	//you can sanitize, cut it down, add images, etc
	policy := bluemonday.StrictPolicy() //remove all html tags
	sanitized := policy.Sanitize(string(blackfriday.MarkdownCommon([]byte(post.Description))))
	excerpt := template.HTML(truncate(sanitized, 300) + "...")
	return excerpt
}

//GetPost returns Post by its ID. Also initializes post author and tags fields.
func GetPost(id string) (*Post, error) {
	post := &Post{}
	collection := db.C("post")
	objectId := bson.ObjectIdHex(id)
	err := collection.Find(bson.M{"_id": objectId}).One(&post)
	logrus.Println("post => ", post.Tags)
	return post, err
}

func GetPostBySlug(slug string) (*Post, error) {
	post := &Post{}
	collection := db.C("post")
	err := collection.Find(bson.M{"slug": slug}).One(&post)
	return post, err
}

//GetPosts returns a slice of posts, order by descending id
func GetPosts() ([]Post, error) {
	var list []Post
	collection := db.C("post")
	err := collection.Find(nil).Sort("-created_at").All(&list)
	logrus.Errorf("list : %s", list)
	return list, err
	//err := db.Select(&list, "SELECT * FROM posts ORDER BY posts.id DESC")
	//return list, err
}

//GetPublishedPosts returns a slice published of posts with their associations
func GetPublishedPosts() ([]Post, error) {
	var list []Post
	collection := db.C("post")
	err := collection.Find(nil).All(&list)
	return list, err
	//var list []Post
	//err := db.Select(&list, "SELECT * FROM posts WHERE published=$1 ORDER BY posts.id DESC", true)
	//if err != nil {
	//	return list, err
	//}
	//if err := fillPostsAssociations(list); err != nil {
	//	return list, err
	//}
	//return list, err
}

//
////GetRecentPosts returns a slice of published posts
func GetRecentPosts() ([]Post, error) {
	var list []Post
	collection := db.C("post")
	err := collection.Find(nil).All(&list)
	//err := db.Select(&list, "SELECT id, name FROM posts WHERE published=$1 ORDER BY id DESC LIMIT 7", true)
	return list, err
}

//
//GetPostMonths returns a slice of distinct months extracted from posts creation dates
func GetPostMonths() ([]Post, error) {
	var list []Post
	post := Post{}
	collection := db.C("post")
	err := collection.Find(bson.M{"slug": "apache-nginx-log"}).One(&post)
	return list, err
	//var list []Post
	//err := db.Select(&list, "SELECT DISTINCT date_trunc('month', created_at) as created_at FROM posts WHERE published=$1 ORDER BY created_at DESC", true)
	//return list, err
}

//
////GetPostsByArchive returns a slice of published posts, given creation year and month
//func GetPostsByArchive(year, month int) ([]Post, error) {
//	var list []Post
//	err := db.Select(&list, "SELECT * FROM posts WHERE published=$1 AND date_part('year', created_at)=$2 AND date_part('month', created_at)=$3 ORDER BY created_at DESC", true, year, month)
//	if err != nil {
//		return list, err
//	}
//	if err := fillPostsAssociations(list); err != nil {
//		return list, err
//	}
//	return list, err
//}
//
//GetPostsByTag returns a slice of published posts associated with tag name
func GetPostsByTag(name string) ([]Post, error) {
	var list []Post
	//err := db.Select(&list, "SELECT * FROM posts WHERE published=$1 AND EXISTS (SELECT null FROM poststags WHERE poststags.post_id=posts.id AND poststags.tag_name=$2) ORDER BY created_at DESC", true, name)
	//if err != nil {
	//	return list, err
	//}
	//if err := fillPostsAssociations(list); err != nil {
	//	return list, err
	//}
	return list, nil
}

func GetPostsByCategory(name string) ([]Post, error) {
	var list []Post
	collection := db.C("post")
	err := collection.Find(bson.M{"category": name}).Sort("-_id").All(&list)
	//logrus.Errorf("tag list : %s", list)
	return list, err
}

////fillPostsAssociations initialises post associations, given post slice
//func fillPostsAssociations(list []Post) error {
//	for i := range list {
//		err := db.Get(&list[i].Author, "SELECT id,name FROM users WHERE id=$1", list[i].UserID)
//		if err != nil {
//			return err
//		}
//		err = db.Select(&list[i].Tags, "SELECT name FROM tags WHERE EXISTS (SELECT null FROM poststags WHERE post_id=$1 AND tag_name=tags.name)", list[i].ID)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
