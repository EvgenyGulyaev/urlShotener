package main

import (
	"fmt"
	"log"
	"urlShortener/internal/config"
	"urlShortener/internal/http/routes"
	"urlShortener/internal/store"
	"urlShortener/pkg/server"
)

func main() {
	c := config.LoadConfig()
	store.InitStore()

	getRoutes := map[string]server.Get{
		"/url":      {Callback: routes.GetUrl},
		"/api/urls": {Callback: routes.GetUrls},
	}
	postRoutes := map[string]server.Post{
		"/url":      {Callback: routes.PostUrl},
		"/api/urls": {Callback: routes.PostUrl},
	}
	deleteRoutes := map[string]server.Delete{
		"/api/urls":     {Callback: routes.DeleteUrl},
		"/api/urls/all": {Callback: routes.ClearUrls},
	}

	s := server.GetServer(fmt.Sprintf(":%s", c.Env["PORT"]), getRoutes, postRoutes, deleteRoutes)
	err := s.StartHandle()
	if err != nil {
		log.Print(err)
	}
}
