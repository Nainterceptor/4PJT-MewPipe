package entities

import (
	"encoding/base64"
	"errors"
	"supinfo/mewpipe/configs"
	"time"

	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const tokenExpiration = 3600

func getUserCollection() *mgo.Collection {
	return configs.MongoDB.C("users")
}

type name struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	NickName  string `json:"nickname"`
}

type UserToken struct {
	Token    bson.ObjectId `json:"token"`
	ExpireAt time.Time     `json:"expireAt"`
}

type User struct {
	Id             bson.ObjectId `json:"id" bson:"_id,omitempty"`
	CreatedAt      time.Time     `json:"createdAt" bson:"createdAt"`
	Name           name          `json:"name" bson:",omitempty"`
	Email          string        `json:"email" bson:",omitempty"`
	Roles          []string      `json:"roles,omitempty" bson",omitempty"`
	Password       string        `json:"password,omitempty" bson:"-"`
	HashedPassword string        `json:"-" bson:",omitempty"`
	UserTokens     []UserToken   `json:"-" bson:",omitempty"`
}

func UserNew() *User {
	user := new(User)
	user.Id = bson.NewObjectId()
	user.CreatedAt = time.Now()

	return user
}

func userTokenNew() *UserToken {
	uToken := new(UserToken)
	currentTime := time.Now()
	uToken.Token = bson.NewObjectIdWithTime(currentTime)
	uToken.ExpireAt = currentTime.Add(time.Second * tokenExpiration)

	return uToken
}

func UserNewFromId(oid bson.ObjectId) *User {
	user := new(User)
	user.Id = oid
	user.CreatedAt = time.Now()

	return user
}

func UserFromId(oid bson.ObjectId) (*User, error) {
	user := new(User)
	err := getUserCollection().FindId(oid).One(&user)
	if err != nil {
		user = UserNewFromId(oid)
	}

	return user, err
}

func UserList(bson bson.M, start int, number int) ([]User, error) {
	users := make([]User, number)

	err := getUserCollection().Find(bson).Skip(start).Limit(number).All(&users)

	return users, err
}

func UserFromCredentials(email string, password string) (*User, error) {
	user := new(User)

	if err := getUserCollection().Find(bson.M{"email": email}).One(&user); err != nil {
		return new(User), err
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return new(User), err
	}
	return user, nil
}

func UserFromToken(token string) (*User, error) {
	user := new(User)
	tokenDecoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return new(User), err
	}
	if err := getUserCollection().Find(bson.M{"usertokens.token": bson.ObjectId(tokenDecoded)}).One(&user); err != nil {
		return new(User), err
	}
	user.Clean()
	if !user.hasToken(string(tokenDecoded)) {
		return new(User), errors.New("Token expired")
	}

	return user, nil
}

func (u *User) Normalize() {
	u.Email = govalidator.Trim(u.Email, "")
	u.Name.FirstName = govalidator.Trim(u.Name.FirstName, "")
	u.Name.LastName = govalidator.Trim(u.Name.LastName, "")
	u.Name.NickName = govalidator.Trim(u.Name.NickName, "")
}

func (u *User) Validate() error {
	u.Normalize()
	if u.Email == "" {
		return errors.New("`email` is empty")
	}
	if u.Name.NickName == "" {
		return errors.New("`nickname` is empty")
	}
	if !govalidator.IsEmail(u.Email) {
		return errors.New("`email` must be valid")
	}
	return nil
}

func (u *User) hasToken(token string) bool {
	oid := bson.ObjectId(token)
	found := false
	for _, token := range u.UserTokens {
		if token.Token == oid {
			found = true
		}
	}
	return found
}

func (u *User) Clean() {
	u.Password = ""
	change := false
	for i, token := range u.UserTokens {
		if token.ExpireAt.Before(time.Now()) {
			u.UserTokens = append(u.UserTokens[:i], u.UserTokens[i+1:]...)
			change = true
		}
	}
	if change {
		u.Update()
	}
}

func (u *User) hashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	u.HashedPassword = string(hashedPassword[:])

	return err
}

func (u *User) HasRole(search string) bool {
	set := make(map[string]struct{}, len(u.Roles))
	for _, s := range u.Roles {
		set[s] = struct{}{}
	}

	_, ok := set[search]
	return ok
}

func (u *User) Insert() error {
	defer u.Clean()

	if u.Password == "" {
		return errors.New("`password` is empty")
	}
	u.hashPassword()
	if err := getUserCollection().Insert(&u); err != nil {
		return err
	}
	return nil
}

func (u *User) Update() error {
	defer u.Clean()

	if u.Password != "" {
		u.hashPassword()
	}

	if err := getUserCollection().UpdateId(u.Id, &u); err != nil {
		return err
	}
	mediaCol, _ := u.GetMedia()
	for _, media := range mediaCol {
		media.Publisher.Id = u.Id
		media.Publisher.Email = u.Email
		media.Publisher.Name = u.Name
		media.Update()
	}

	return nil
}

func (u *User) Delete() error {
	if err := getUserCollection().RemoveId(u.Id); err != nil {
		return err
	}
	mediaCol, _ := u.GetMedia()
	for _, media := range mediaCol {
		media.Delete()
	}
	return nil
}

func (u *User) GetMedia() ([]Media, error) {
	var media []Media
	err := getMediaCollection().Find(bson.M{"publisher._id": u.Id}).All(&media)
	return media, err
}

func (u *User) TokenNew() (*UserToken, error) {
	newToken := userTokenNew()
	u.UserTokens = append(u.UserTokens, *newToken)
	if err := u.Update(); err != nil {
		return new(UserToken), err
	}
	return newToken, nil
}
