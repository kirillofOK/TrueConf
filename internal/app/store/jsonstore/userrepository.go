package jsonstore

import (
	"github.com/kirillofOK/TrueConf/internal/app/model"
	"github.com/kirillofOK/TrueConf/internal/app/store"
)

type UserRepository struct {
	//store *Store
	increment int
	users     map[int]*model.User
}

func (r *UserRepository) Create(u *model.User) error {

	if err := u.Vaidate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return nil
	}

	u.ID = len(r.users) + 1
	r.users[u.ID] = u
	return nil
}

func (r *UserRepository) Find(id int) (*model.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}

func (r *UserRepository) Delete(u *model.User) error {
	// Access rights should be checked here
	delete(r.users, u.ID)
	return nil
}
