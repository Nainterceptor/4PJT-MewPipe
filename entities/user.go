package entities

import(
    "supinfo/mewpipe/configs"
    "gopkg.in/mgo.v2/bson"
    "golang.org/x/crypto/bcrypt"
    "errors"
)

var userCollection = configs.MongoDB.C("users")

type name struct {
    FirstName   string  `json:"firstname"`
    LastName    string  `json:"lastname"`
    NickName    string  `json:"nickname"`
}

type UserToken struct {
    Token   string `json:"token"`
}

type User struct {
    Id              bson.ObjectId   `json:"id" bson:"_id,omitempty"`
    Name            name            `json:"name" bson:",omitempty"`
    Email           string          `json:"email" bson:",omitempty"`
    Roles           []string        `json:"-" bson",omitempty"`
    Password        string          `json:"password,omitempty" bson:",omitempty"`
    HashedPassword  string          `json:"-" bson:",omitempty"`
    UserTokens      []UserToken     `json:"-" bson:",omitempty"`
}

func UserNew() *User {
    user := new(User)
    user.Id = bson.NewObjectId()
    return user
}

func UserNewFromId(oid bson.ObjectId) *User {
    user := new(User)
    user.Id = bson.NewObjectId()
    return user
}

func UserFromId(oid bson.ObjectId) (*User, error) {
    user := new(User)
    err := userCollection.FindId(oid).One(user);
    if err != nil {
        user = UserNewFromId(oid)
    }

    return user, err
}

func (u *User) Validate() error {
    if u.Email == "" {
        return errors.New("`email` is empty")
    }
    if u.Name.NickName == "" {
        return errors.New("`nickname` is empty")
    }
    return nil
}

func (u *User) clean() {
    u.Password = ""
}

func (u *User) hashPassword() error {
    defer u.clean()

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
    u.HashedPassword = string(hashedPassword[:])

    return err
}

func (u *User) Insert() error {
    if u.Password == "" {
        return errors.New("`password` is empty")
    }
    u.hashPassword()
    if err := userCollection.Insert(&u); err != nil {
        return err
    }
    return nil
}

func (u *User) Update() error {
    if u.Password != "" {
        u.hashPassword()
    }

    if err := userCollection.UpdateId(u.Id, &u); err != nil {
        return err
    }
    return nil
}
