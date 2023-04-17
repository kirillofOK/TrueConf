package store

import (
	"github.com/kirillofOK/TrueConf/internal/app/model"
)

type UserRepository interface {
	Create(*model.User) error
	SearchUsers() *map[string]*model.User
	Get(id string) (*model.User, error)
	Update(id string, display_name string) error
	Delete(id string) error
}
