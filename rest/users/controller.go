package users

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"supinfo/mewpipe/entities"
	"time"

	"strconv"

	"supinfo/mewpipe/configs"

	"github.com/Nainterceptor/go-restful"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	goth.UseProviders(
		twitter.New(*configs.TWITTER_KEY, *configs.TWITTER_SECRET, "http://localhost:1337/rest/users/login/twitter/callback"),
	)
}

func userCreate(request *restful.Request, response *restful.Response) {

	usr := entities.UserNew()

	if err := request.ReadEntity(&usr); err != nil {
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	usr.Roles = usr.Roles[:0]
	if err := usr.Validate(); err != nil {
		response.WriteError(http.StatusNotAcceptable, err)
		return
	}
	if usr.Password == "" {
		response.WriteErrorString(http.StatusNotAcceptable, "`password` is empty")
	}

	if err := usr.Insert(); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.WriteEntity(usr)
}

func usersGet(request *restful.Request, response *restful.Response) {
	start, err := strconv.Atoi(request.QueryParameter("start"))
	if err != nil {
		start = 0
	}
	limit, err := strconv.Atoi(request.QueryParameter("limit"))
	if err != nil {
		limit = 25
	}
	if limit < 1 || limit > 100 {
		limit = 25
	}
	users, err := entities.UserList(bson.M{}, start, limit)

	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.WriteEntity(users)
}

func userUpdate(request *restful.Request, response *restful.Response) {

	id := request.PathParameter("user-id")
	if !bson.IsObjectIdHex(id) {
		response.WriteErrorString(http.StatusBadRequest, "Path must contain an Object ID")
		return
	}

	usr, err := entities.UserFromId(bson.ObjectIdHex(id))
	if err != nil {
		response.WriteError(http.StatusNotFound, err)
		return
	}
	oldUsr := *usr
	if err := request.ReadEntity(&usr); err != nil {
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	//Persist some informations
	usr.Id = oldUsr.Id
	usr.Roles = oldUsr.Roles

	if err := usr.Validate(); err != nil {
		response.WriteError(http.StatusNotAcceptable, err)
		return
	}

	if err := usr.Update(); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	response.WriteEntity(usr)
}

func userDelete(request *restful.Request, response *restful.Response) {

	id := request.PathParameter("user-id")
	if !bson.IsObjectIdHex(id) {
		response.WriteErrorString(http.StatusBadRequest, "Path must contain an Object ID")
		return
	}

	//userNew because find query is useless
	usr := entities.UserNewFromId(bson.ObjectIdHex(id))

	if err := usr.Delete(); err != nil {
		response.WriteError(http.StatusNotFound, err)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

func userGet(request *restful.Request, response *restful.Response) {

	id := request.PathParameter("user-id")
	if !bson.IsObjectIdHex(id) {
		response.WriteErrorString(http.StatusBadRequest, "Path must contain an Object ID")
		return
	}

	usr, err := entities.UserFromId(bson.ObjectIdHex(id))
	if err != nil {
		response.WriteError(http.StatusNotFound, err)
		return
	}

	response.WriteEntity(usr)
}

func userLogin(request *restful.Request, response *restful.Response) {

	form := entities.UserNew()

	if err := request.ReadEntity(&form); err != nil {
		response.WriteError(http.StatusBadRequest, err)
		return
	}
	usr, err := entities.UserFromCredentials(form.Email, form.Password)
	if err != nil {
		response.WriteError(http.StatusNotFound, err)
		return
	}
	token, err := usr.TokenNew()
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	type lambdaReturn struct {
		User     entities.User
		Token    string
		ExpireAt time.Time
	}

	toReturn := new(lambdaReturn)
	toReturn.User = *usr
	toReturn.ExpireAt = token.ExpireAt
	toReturn.Token = base64.StdEncoding.EncodeToString([]byte(token.Token))

	response.WriteEntity(toReturn)
}

func getAuthURL(request *restful.Request, response *restful.Response) (string, error) {

	providerName := request.PathParameter("provider")
	res := response.ResponseWriter
	req := request.Request

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}
	sess, err := provider.BeginAuth(gothic.GetState(req))
	if err != nil {
		return "", err
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}

	session, _ := gothic.Store.Get(req, gothic.SessionName)
	session.Values[gothic.SessionName] = sess.Marshal()
	err = session.Save(req, res)
	if err != nil {
		return "", err
	}

	return url, err
}

func userThirdPartyLogin(request *restful.Request, response *restful.Response) {
	url, err := getAuthURL(request, response)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(response, err)
		return
	}

	http.Redirect(response.ResponseWriter, request.Request, url, http.StatusTemporaryRedirect)
}

func userThirdPartyLoginCallback(request *restful.Request, response *restful.Response) {
	providerName := request.PathParameter("provider")
	req := request.Request

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		response.WriteError(http.StatusNotFound, err)
		return
	}

	session, _ := gothic.Store.Get(req, gothic.SessionName)

	if session.Values[gothic.SessionName] == nil {
		response.WriteErrorString(http.StatusNotFound, "could not find a matching session for this request")
		return
	}

	sess, err := provider.UnmarshalSession(session.Values[gothic.SessionName].(string))
	if err != nil {
		response.WriteError(http.StatusNotFound, err)
		return
	}

	_, err = sess.Authorize(provider, req.URL.Query())
	if err != nil {
		response.WriteError(http.StatusNotFound, err)
		return
	}

	user, err := provider.FetchUser(sess)
	if err != nil {
		response.WriteError(http.StatusNotFound, err)
		return
	}

	usr, err := entities.UserFromTwitterUserID(user.UserID)
	if err != nil {
		usr := entities.UserNew()
		usr.Name.NickName = user.NickName
		usr.Twitter.UserId = user.UserID

		if err := usr.Insert(); err != nil {
			response.WriteError(http.StatusInternalServerError, err)
			fmt.Println("1")

			return
		}
	}

	token, err := usr.TokenNew()
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		fmt.Println("1")

		return
	}

	type lambdaReturn struct {
		User     entities.User
		Token    string
		ExpireAt time.Time
	}

	toReturn := new(lambdaReturn)
	toReturn.User = *usr
	toReturn.ExpireAt = token.ExpireAt
	toReturn.Token = base64.StdEncoding.EncodeToString([]byte(token.Token))

	response.WriteEntity(toReturn)
}

func userRefreshToken(request *restful.Request, response *restful.Response) {

	usr := request.Attribute("user").(*entities.User)
	//Not useful to remove current token, because we've an automatic clean for older tokens.
	token, err := usr.TokenNew()
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	type tokenBack struct {
		Token    string
		ExpireAt time.Time
	}

	theToken := new(tokenBack)
	theToken.ExpireAt = token.ExpireAt
	theToken.Token = base64.StdEncoding.EncodeToString([]byte(token.Token))

	response.WriteEntity(theToken)
}
