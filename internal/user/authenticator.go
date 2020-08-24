package user

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
	"time"

	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
	"github.com/Odar/capital-lurker/pkg/app/repositories"
	"github.com/labstack/echo/v4"
)

func New(repo repositories.AuthenticatorRepo) *authenticator {
	return &authenticator{
		repo: repo,
	}
}

const secret = "Please, change me!"

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
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
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

func (a *authenticator) signUp(email, password, firstName, lastName string, birthDate time.Time) (*models.User, error) {
	//add a table with emails to validate and another method to move user from tmp table to user table
	model, err := a.repo.AddUser(email, password, firstName, lastName, birthDate)
	return model, err
}
