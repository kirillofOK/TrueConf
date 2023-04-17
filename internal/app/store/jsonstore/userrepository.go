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

func (r *UserRepository) Create(u *model.User) error {

	if err := u.Validate(); err != nil {
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

func (r *UserRepository) SearchUsers() *map[string]*model.User {

	r.store.Open()
	return &r.users
}

func (r *UserRepository) Get(id string) (*model.User, error) {
	r.store.Open()
	u, ok := r.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}

func (r *UserRepository) Delete(id string) error {
	r.store.Open()
	u, ok := r.users[id]
	if !ok {
		return store.ErrRecordNotFound
	}
	// Access rights should be checked here
	delete(r.users, strconv.Itoa(u.ID))
	r.store.Save()
	return nil
}

func (r *UserRepository) Update(id string, display_name string) error {
	r.store.Open()
	u, ok := r.users[id]
	if !ok {
		return store.ErrRecordNotFound
	}

	u.DisplayName = display_name

	r.users[id] = u

	r.store.Save()

	return nil
}
