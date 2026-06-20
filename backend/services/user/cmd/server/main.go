package main

import (
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/cmd/server/internal/handler"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/cmd/server/internal/repository"
	"github.com/mak-magz/myconfed-microsvc/backend/services/user/cmd/server/internal/service"
)

func main() {
	// wiring: repository -> service -> handler
	repo := repository.NewRepository()
	svc := service.NewService(repo)
	hnd := handler.NewHandler(svc)
	hnd.GetUser("123")
}
