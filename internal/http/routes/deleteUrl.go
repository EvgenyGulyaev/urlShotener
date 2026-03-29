package routes

import (
	"log"
	"net/http"
	"urlShortener/internal/store"
	"urlShortener/pkg/server/callback"

	"github.com/go-www/silverlining"
)

func DeleteUrl(ctx *silverlining.Context) {
	u := store.GetUrlRepository()

	short, shortErr := ctx.GetQueryParamString("short")
	if shortErr == nil && short != "" {
		if err := u.DeleteByShort(short); err != nil {
			callback.GetError(ctx, &callback.Error{Message: err.Error(), Status: http.StatusInternalServerError})
			return
		}

		writeDeleteResponse(ctx, short)
		return
	}

	original, urlErr := ctx.GetQueryParamString("url")
	if urlErr == nil && original != "" {
		if err := u.DeleteLink(original); err != nil {
			callback.GetError(ctx, &callback.Error{Message: err.Error(), Status: http.StatusInternalServerError})
			return
		}

		writeDeleteResponse(ctx, original)
		return
	}

	callback.GetError(ctx, &callback.Error{
		Message: "short or url query param is required",
		Status:  http.StatusBadRequest,
	})
}

func ClearUrls(ctx *silverlining.Context) {
	u := store.GetUrlRepository()

	if err := u.ClearAll(); err != nil {
		callback.GetError(ctx, &callback.Error{Message: err.Error(), Status: http.StatusInternalServerError})
		return
	}

	err := ctx.WriteJSON(http.StatusOK, map[string]string{"status": "cleared"})
	if err != nil {
		log.Print(err)
	}
}

func writeDeleteResponse(ctx *silverlining.Context, value string) {
	err := ctx.WriteJSON(http.StatusOK, map[string]string{
		"status":  "deleted",
		"deleted": value,
	})
	if err != nil {
		log.Print(err)
	}
}
