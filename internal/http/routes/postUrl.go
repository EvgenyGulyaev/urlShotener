package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"urlShortener/internal/store"
	"urlShortener/pkg/server/callback"

	"github.com/go-www/silverlining"
)

type bodyPostUrl struct {
	Url string `json:"url"`
}

func PostUrl(ctx *silverlining.Context, body []byte) {
	var req bodyPostUrl
	err := json.Unmarshal(body, &req)
	if err != nil {
		callback.GetError(ctx, &callback.Error{Message: err.Error(), Status: http.StatusInternalServerError})
		return
	}

	u := store.GetUrlRepository()
	url, err := u.Create(req.Url)
	if err != nil {
		callback.GetError(ctx, &callback.Error{Message: err.Error(), Status: http.StatusInternalServerError})
		return
	}

	err = ctx.WriteJSON(http.StatusOK, url)
	if err != nil {
		log.Print(err)
	}
}
