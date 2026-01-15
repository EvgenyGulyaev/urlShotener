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
		"/url": {Callback: routes.GetUrl},
	}
	postRoutes := map[string]server.Post{
		"/url": {Callback: routes.PostUrl},
	}

	s := server.GetServer(fmt.Sprintf(":%s", c.Env["PORT"]), getRoutes, postRoutes)
	err := s.StartHandle()
	if err != nil {
		log.Print(err)
	}
}
