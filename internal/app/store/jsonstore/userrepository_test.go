package jsonstore_test

import (
	"testing"

	"github.com/kirillofOK/TrueConf/internal/app/model"
	"github.com/kirillofOK/TrueConf/internal/app/store/jsonstore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	s := jsonstore.New("/Users/olegkirillov/go/src/TrueConf/users.json")
	u := model.TestUser(t)

	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)

}
