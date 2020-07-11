package receiver

import (
	"encoding/json"
	"net/http"

	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/labstack/echo/v4"
)

func New() *receiver {
	return &receiver{}
}

type receiver struct {
}

func (r *receiver) Reverse(ctx echo.Context) error {
	var request api.ReverseRequest
	err := ctx.Bind(&request)
	if err != nil {
		return ctx.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	n := 0
	inRunes := make([]rune, len(request.Word))
	for _, r := range request.Word {
		inRunes[n] = r
		n++
	}
	inRunes = inRunes[0:n]

	// Reverse
	for i := 0; i < n/2; i++ {
		inRunes[i], inRunes[n-1-i] = inRunes[n-1-i], inRunes[i]
	}
	// Convert back to UTF-8.
	reversed := string(inRunes)

	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(api.ReverseResponse{
		Word: reversed,
	})
}
