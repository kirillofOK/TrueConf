package jsonstore

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"

	"github.com/kirillofOK/TrueConf/internal/app/model"
	"github.com/kirillofOK/TrueConf/internal/app/store"
)

type Store struct {
	fileURL        string
	userRepository *UserRepository
}

func New(f string) *Store {
	return &Store{
		fileURL: f,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[string]*model.User),
	}

	return s.userRepository
}

func (s *Store) Open() error {
	f, err := ioutil.ReadFile(s.fileURL)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(f, &s.userRepository.users); err != nil {
		return err
	}
	return nil
}

func (s *Store) Save() error {
	b, err := json.Marshal(&s.userRepository.users)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(s.fileURL, b, fs.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
