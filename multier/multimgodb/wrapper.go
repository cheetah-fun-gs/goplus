package multimgodb

import (
	"github.com/globalsign/mgo"
)

// Clone ...
func Clone() *mgo.Session {
	s, _ := CloneN(d)
	return s
}

// CloneN ...
func CloneN(name string) (*mgo.Session, error) {
	db, err := RetrieveN(name)
	if err != nil {
		return nil, err
	}
	return db.Session.Clone(), nil
}

// MustCloneN ...
func MustCloneN(name string) *mgo.Session {
	s, err := CloneN(d)
	if err != nil {
		panic(err)
	}
	return s
}
