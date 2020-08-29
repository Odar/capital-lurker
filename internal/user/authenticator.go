package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-vk-api/vk"
	"golang.org/x/oauth2"
	vkAuth "golang.org/x/oauth2/vk"
	"net/http"
	"time"

	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
	"github.com/Odar/capital-lurker/pkg/app/repositories"
	"github.com/labstack/echo/v4"
)

var (
	vkOauthConfig = &oauth2.Config{
		ClientID:     "7562746",                                 //os.Getenv("CLIENT_ID"),
		ClientSecret: "rqozJXDcBPrygXS181xr",                    //os.Getenv("CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8888/tmppageforvkoauth", //os.Getenv("REDIRECT_URL"),
		Scopes:       []string{},
		Endpoint:     vkAuth.Endpoint,
	}
	randomState = "haha"
)

const secret = "Please, change me!"

type vkUser struct {
	ID        uint64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type vkUserInfo struct {
	ID        uint64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"bdate"`
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

	t, err := a.createToken(id)
	if err != nil {
		return err
	}

	err = a.repo.UpdateTokenWithID(id, t)
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
	id := uint64(claims["id"].(float64))
	err := a.repo.InvalidateToken(id)
	if err != nil {
		return err
	}

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
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": claims["id"],
	})
}

func (a *authenticator) LoginVkCheckRegistration(ctx echo.Context) error {
	//validate with db
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

	id, err := a.getUserVkID(client)
	if err != nil {
		return err
	}

	isRegistrated, err := a.repo.CheckRegistration(id)
	if err != nil {
		return err
	}

	if !isRegistrated {
		return a.registrationVk(ctx, client)
	}

	return a.loginVk(ctx, id)
}

func (a *authenticator) registrationVk(ctx echo.Context, api *vk.Client) error {
	vkUser, err := a.getUserVkInfo(api)
	if err != nil {
		return err
	}

	email, err := a.getEmailForVkRegistration(ctx)
	if err != nil {
		return err
	}

	_, err = a.repo.AddUser(email, "", vkUser.FirstName, vkUser.LastName, time.Time{}, vkUser.ID) //add nullable to db, get birthDate from string
	if err != nil {
		return err
	}

	return a.loginVk(ctx, vkUser.ID)
}

func (a *authenticator) getEmailForVkRegistration(ctx echo.Context) (string, error) {
	return "", nil
}

func (a *authenticator) loginVk(ctx echo.Context, vkID uint64) error {
	id, err := a.repo.GetIDForVk(vkID)
	if err != nil {
		return err
	}

	token, err := a.createToken(id)
	if err != nil {
		return err
	}

	err = a.repo.UpdateTokenWithID(id, token)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func (a *authenticator) getUserVkID(api *vk.Client) (uint64, error) {
	var users []vkUser
	err := api.CallMethod("users.get", vk.RequestParams{
		"v":      "5.122",
		"fields": "",
	}, &users)
	if err != nil {
		return 0, err
	}
	if len(users) < 1 {
		return 0, errors.New("returned blank list of users")
	}

	return users[0].ID, nil
}

func (a *authenticator) getUserVkInfo(api *vk.Client) (vkUserInfo, error) {
	var users []vkUserInfo
	err := api.CallMethod("users.get", vk.RequestParams{
		"v":      "5.122",
		"fields": "bdate",
	}, &users)
	if err != nil {
		return vkUserInfo{}, err
	}
	if len(users) < 1 {
		return vkUserInfo{}, errors.New("returned blank list of users")
	}

	return users[0], nil
}

func (a *authenticator) LoginVkInitOauth(ctx echo.Context) error {
	url := vkOauthConfig.AuthCodeURL(randomState)
	//storing randomState in db
	return ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (a *authenticator) signUp(email, password, firstName, lastName string, birthDate time.Time) (*models.User, error) {
	model, err := a.repo.AddUser(email, password, firstName, lastName, birthDate, 0)
	return model, err
}

func (a *authenticator) createToken(id uint64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour).Unix() // session duration

	tokenString, err := token.SignedString([]byte(secret))
	err = a.repo.UpdateTokenWithID(id, tokenString)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func (a *authenticator) CheckTokenValidityMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token := ctx.Get("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		id := uint64(claims["id"].(float64))
		result, err := a.repo.CheckToken(id, token.Raw)
		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Message:  "cannot get jwt",
				Internal: err,
			}
		}
		if result == repositories.TokenStatusInvalidToken {
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  "invalid or expired token",
				Internal: err,
			}
		}
		if result == repositories.TokenStatusLoggedOut {
			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  "logged out",
				Internal: err,
			}
		}

		return next(ctx)
	}
}
