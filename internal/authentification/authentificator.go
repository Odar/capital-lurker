package authentification

import (
	"encoding/json"
	"fmt"
	vk "github.com/go-vk-api/vk"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	vkAuth "golang.org/x/oauth2/vk"
	"net/http"
)

type authentificator struct {
}

func New() *authentificator {
	return &authentificator{}
}

var (
	vkOauthConfig = &oauth2.Config{
		ClientID:     "7562746",                    //os.Getenv("CLIENT_ID"),
		ClientSecret: "rqozJXDcBPrygXS181xr",       //os.Getenv("CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8888/home", //os.Getenv("REDIRECT_URL"),
		Scopes:       []string{},
		Endpoint:     vkAuth.Endpoint,
	}
	randomState = "haha"
)

func (a *authentificator) GetInfoFromVK(ctx echo.Context) error {
	if ctx.FormValue("state") != randomState {
		return ctx.String(http.StatusInternalServerError, "state is not valid")
	}
	token, err := vkOauthConfig.Exchange(oauth2.NoContext, ctx.FormValue("code"))
	if err != nil {
		return ctx.String(http.StatusInternalServerError, fmt.Sprintf("can not get token %s\n", err.Error()))
	}

	client, err := vk.NewClientWithOptions(vk.WithToken(token.AccessToken))
	if err != nil {
		return ctx.String(http.StatusInternalServerError, fmt.Sprintf("can not create vk api client %s\n", err.Error()))
	}

	user := getCurrentUser(client)
	return json.NewEncoder(ctx.Response()).Encode(user)
}

func getCurrentUser(api *vk.Client) User {
	var users []User

	api.CallMethod("users.get", vk.RequestParams{
		"fields": "photo_400_orig,city",
	}, &users)

	return users[0]
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Photo     string `json:"photo_400_orig"`
	City      City   `json:"city"`
}

type City struct {
	Title string `json:"title"`
}

func (a *authentificator) Login(ctx echo.Context) error {
	url := vkOauthConfig.AuthCodeURL(randomState)
	err := ctx.Redirect(http.StatusTemporaryRedirect, url)
	return err
}
