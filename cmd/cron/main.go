package main

import (
	"time"
	"urlShortener/internal/config"
	"urlShortener/internal/store"
)

func main() {
	config.LoadConfig()
	store.InitStore()

	repo := store.GetUrlRepository()
	urlCleaner := repo.StartAutoCleanup(time.Hour)
	defer urlCleaner.Stop()
}
