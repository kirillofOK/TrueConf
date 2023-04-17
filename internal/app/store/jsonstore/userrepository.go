package jsonstore

import (
	"strconv"

	"github.com/kirillofOK/TrueConf/internal/app/model"
	"github.com/kirillofOK/TrueConf/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[string]*model.User
}

//
func (r *UserRepository) Create(u *model.User) error {

	if err := u.Vaidate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	if err := r.store.Open(); err != nil {
		return err
	}

	u.ID = len(r.users) + 1
	r.users[strconv.Itoa(u.ID)] = u

	if err := r.store.Save(); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Find(id int) (*model.User, error) {
	u, ok := r.users[strconv.Itoa(id)]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}

func (r *UserRepository) Delete(u *model.User) error {
	// Access rights should be checked here
	delete(r.users, strconv.Itoa(u.ID))
	return nil
}

func (r *UserRepository) Update(id int, display_name string) error {
	u, ok := r.users[strconv.Itoa(id)]
	if !ok {
		return store.ErrRecordNotFound
	}

	u.DisplayName = display_name

	return nil
}
