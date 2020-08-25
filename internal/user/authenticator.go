package user

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-vk-api/vk"
	"golang.org/x/oauth2"
	vkAuth "golang.org/x/oauth2/vk"
	"net/http"
	"strconv"
	"time"

	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
	"github.com/Odar/capital-lurker/pkg/app/repositories"
	"github.com/labstack/echo/v4"
)

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

const secret = "Please, change me!"

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BDate     string `json:"bdate"`
	Number    tNum   `json:"contacts"`
}

type tNum struct {
	Mobile string `json:"mobile_phone"`
	Home   string `json:"home_phone"`
}

func New(repo repositories.AuthenticatorRepo) *authenticator {
	return &authenticator{
		repo: repo,
	}
}

type authenticator struct {
	repo repositories.AuthenticatorRepo
}

func (a *authenticator) Login(ctx echo.Context) error {
	var request api.SignInRequest
	err := ctx.Bind(&request)
	if err != nil {
		return ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	id, authorized, err := a.repo.CheckAuth(request.Email, request.Password)
	if err != nil {
		return err
	}
	if !authorized {
		return ctx.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
	}
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = strconv.FormatUint(uint64(id), 10)
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func (a *authenticator) Logout(ctx echo.Context) error {
	token := ctx.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Unix(0, 0)
	return ctx.String(http.StatusUnauthorized, "You are no longer authorized")
}

func (a *authenticator) SignUp(ctx echo.Context) error {
	var request api.SignUpRequest
	err := ctx.Bind(&request)
	if err != nil {
		return ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	model, err := a.signUp(request.Email, request.Password, request.FirstName, request.LastName, request.BirthDate)
	if err != nil {
		return ctx.String(http.StatusConflict, err.Error())
	}

	ctx.Response().WriteHeader(http.StatusOK)
	if model == nil {
		return json.NewEncoder(ctx.Response()).Encode(api.SignUpResponse{
			Result: *model,
		})
	}
	return json.NewEncoder(ctx.Response()).Encode(api.SignUpResponse{
		Result: *model,
	})
}

func (a *authenticator) TestPage(ctx echo.Context) error {
	token := ctx.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	return ctx.String(http.StatusOK, claims["id"].(string))
}

func (a *authenticator) SignUpViaVk(ctx echo.Context) error {

	return nil
}

func (a *authenticator) GetInfoFromVK(ctx echo.Context) error {
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
		"v":      "5.122",
		"fields": "bdate,contacts",
	}, &users)

	return users[0]
}

func (a *authenticator) Loginvk(ctx echo.Context) error {
	url := vkOauthConfig.AuthCodeURL(randomState)
	return ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (a *authenticator) signUp(email, password, firstName, lastName string, birthDate time.Time) (*models.User, error) {
	model, err := a.repo.AddUser(email, password, firstName, lastName, birthDate)
	return model, err
}
