package store

import (
	"log"
	"urlShortener/pkg/db"
)

var (
	UrlBucket = []byte("Url")
)

func InitStore() {
	repo := db.GetRepository()
	err := repo.EnsureBuckets([][]byte{UrlBucket})
	if err != nil {
		log.Println(err)
	}
}
