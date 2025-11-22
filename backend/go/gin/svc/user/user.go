package user

import (
	"fmt"

	mdb "github.com/unifuu/hitotose/gin/db/mongo"
	"github.com/unifuu/hitotose/gin/model/user"
	mgo "github.com/unifuu/monggo"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	ByUsername(username string) user.User
	SignIn(username, password string) (user.User, error)
	SignUp(u user.User) error
}

func NewService() Service {
	return &service{}
}

type service struct{}

func (s *service) ByUsername(username string) user.User {
	var u user.User
	filter := bson.D{primitive.E{Key: "username", Value: username}}
	mgo.FindOne(mdb.Users, filter).Decode(&u)
	return u
}

func (s *service) SignIn(username, password string) (user.User, error) {
	u := s.ByUsername(username)
	if len(u.Username) > 0 && bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil {
		return u, nil
	} else {
		return user.User{}, fmt.Errorf("login failed")
	}
}

func (s *service) SignUp(u user.User) error {
	find := s.ByUsername(u.Username)
	if len(find.Username) == 0 {
		return mgo.Insert(mdb.Users, u)
	} else {
		return fmt.Errorf("register failed")
	}
}
