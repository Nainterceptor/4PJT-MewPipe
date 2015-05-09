package rest

import (
	"github.com/emicklei/go-restful"
	"supinfo/mewpipe/entities"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"supinfo/mewpipe/utils"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"time"
	"fmt"
)

func UserRoute() *restful.WebService {
	service := new(restful.WebService)
	service.
	Path("/rest/users").
	Consumes(restful.MIME_JSON).
	Produces(restful.MIME_JSON)

	service.Route(service.POST("/login").To(Connexion)).
		Doc("Login").
		Param(service.BodyParameter("password", "the form password").DataType("string"))
	service.Route(service.POST("").To(CreateUser))
	service.Route(service.GET("").To(GetAllUsers)).
		Doc("get all users")
	service.Route(service.GET("/{user-id}").To(GetUser)).
		Doc("get a user")
	service.Route(service.PUT("/update/{user-id}").To(UpdateUser)).
		Doc("update a single user")
	service.Route(service.DELETE("/delete/{user-id}").To(DeleteUser)).
		Doc("delete a single user")

	return service
}

func NewToken(usr entities.User, mySigningKey []byte) (error) {

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["userid"] = usr.Id
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	tokenString, err := token.SignedString(mySigningKey)

	usrToken := entities.UserToken{ tokenString }

	usr.UserTokens = append(usr.UserTokens, usrToken)

	if err := entities.UserCollection.UpdateId(usr.Id, &usr); err != nil {
		panic(err)
	}

	return err
}

func CheckToken(myToken string) (error) {

	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return token.Header["mypipe"], nil
	})

	if token.Valid {
		fmt.Println("You look nice today")
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			fmt.Println("Timing is everything", err)
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}

	return err
}

func Connexion(request *restful.Request, response *restful.Response) {

	connexion := entities.User{}
	usr := entities.User{}
	errRE := request.ReadEntity(&connexion)

	err := entities.UserCollection.Find(bson.M{"email": connexion.Email}).One(&usr)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.HashedPassword), []byte(connexion.Password))
	if err != nil {
		panic(err)
	}

	if errRE == nil {
		NewToken(usr, []byte("mypipe"))
		if err != nil {
			panic(err)
		}

		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, errRE)
	}
}

func CreateUser(request *restful.Request, response *restful.Response) {

	usr := entities.User{}
	errRE := request.ReadEntity(&usr)
	usr.Id = bson.NewObjectId()
	usr.HashedPassword = utils.Hash(usr.Password)
	usr.Password = ""

	if err := entities.UserCollection.Insert(&usr); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	if errRE == nil {
		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, errRE)
	}
}

func GetAllUsers(request *restful.Request, response *restful.Response) {

	usr := []entities.User{}
	if err := entities.UserCollection.Find(nil).All(&usr); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.WriteEntity(usr)
}

func GetUser(request *restful.Request, response *restful.Response) {

	myToken := request.QueryParameter("token")
	id := request.PathParameter("user-id")

	err := CheckToken(myToken)
	if err != nil {
		panic(err)
	}

	if !bson.IsObjectIdHex(id) {
		response.WriteErrorString(404, "Problem with the id")
		return
	}

	oid := bson.ObjectIdHex(id)
	usr := entities.User{}

	if err := entities.UserCollection.FindId(oid).One(&usr); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.WriteEntity(usr)
}

func UpdateUser(request *restful.Request, response *restful.Response) {

	id := request.PathParameter("user-id")

	if !bson.IsObjectIdHex(id) {
		response.WriteErrorString(404, "Problem with the id")
		return
	}

	oid := bson.ObjectIdHex(id)
	usr := entities.User{}
	errRE := request.ReadEntity(&usr)

	usr.HashedPassword = utils.Hash(usr.Password)
	usr.Password = ""

	if _, err := entities.UserCollection.UpsertId(oid, &usr); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	if errRE == nil {
		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, errRE)
	}
}

func DeleteUser(request *restful.Request, response *restful.Response) {

	id := request.PathParameter("user-id")

	if !bson.IsObjectIdHex(id) {
		response.WriteErrorString(404, "Problem with the id")
		return
	}

	oid := bson.ObjectIdHex(id)
	usr := entities.User{}
	errRE := request.ReadEntity(&usr)

	if err := entities.UserCollection.RemoveId(oid); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	if errRE == nil {
		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, errRE)
	}
}
