package routes

import (
	"log"
	"net/http"
	"urlShortener/internal/store"
	"urlShortener/pkg/server/callback"

	"github.com/go-www/silverlining"
)

func GetUrl(ctx *silverlining.Context) {
	short, err := ctx.GetQueryParamString("url")
	if err != nil {
		callback.GetError(ctx, &callback.Error{Message: err.Error(), Status: http.StatusInternalServerError})
		return
	}

	u := store.GetUrlRepository()
	url, err := u.FindByShort(short)

	if err != nil {
		callback.GetError(ctx, &callback.Error{Message: err.Error(), Status: http.StatusInternalServerError})
		return
	}

	err = ctx.WriteJSON(http.StatusOK, url)
	if err != nil {
		log.Print(err)
	}
}
