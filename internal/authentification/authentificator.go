package authentification

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"
	"github.com/urShadow/go-vk-api"
)

type authentificator struct {
}

func New() *authentificator {
	return &authentificator{}
}

var (
	vkOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       []string{},
		Endpoint:     vk.Endpoint,
	}
	randomState = "haha"
)

func (a *authentificator) GetInfoFromVK(ctx echo.Context) error {
	if ctx.FormValue("state") != randomState {
		return ctx.String(http.StatusInternalServerError, "state is not valid")
	}
	token, err := vkOauthConfig.Exchange(oauth2.NoContext, ctx.FormValue("code"))
	if err != nil {
		return ctx.String(http.StatusInternalServerError, fmt.Sprintf("could not get token %s\n", err.Error()))
	}

	client, err := go_vk_api.
}

func (a *authentificator) Login(ctx echo.Context) error {
	url := vkOauthConfig.AuthCodeURL(randomState)
	ctx.Redirect(http.StatusOK, url)
	return nil
}
