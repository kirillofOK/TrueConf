package store

import (
	"github.com/kirillofOK/TrueConf/internal/app/model"
)

type UserRepository interface {
	Create(*model.User) error
	//Delete(*model.User) error
	//Get()
	//Find(id int) (*model.User, error)
	//Update(id int, display_name string) error
}
