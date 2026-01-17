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
		getUrls(ctx)
		return
	}

	u := store.GetUrlRepository()
	url, err := u.FindByShort(short)
	if err != nil {
		callback.GetError(ctx, &callback.Error{Message: err.Error(), Status: http.StatusInternalServerError})
		return
	}

	go func() {
		err = u.IncrementClicks(short)
		if err != nil {
			log.Print(err)
			return
		}
	}()

	ctx.Redirect(http.StatusFound, url.Original)
}

func getUrls(ctx *silverlining.Context) {
	u := store.GetUrlRepository()
	urls, err := u.ListAll()
	if err != nil {
		callback.GetError(ctx, &callback.Error{Message: err.Error(), Status: http.StatusInternalServerError})
	}
	err = ctx.WriteJSON(http.StatusOK, urls)
	if err != nil {
		log.Print(err)
	}
	return
}
