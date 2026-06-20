package repository

import "fmt"

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetUser(id string) {
	fmt.Println("repository: GetUser", id)
}
