package userinfo

import (
	"os/user"
	"testing"

	"github.com/milosgajdos83/servpeek/utils/group"
	"github.com/stretchr/testify/assert"
)

func Test_FakeProvider_LookupUser_Exists(t *testing.T) {
	expected := &user.User{
		Username: "user",
	}
	p := MakeFakeProvider(501, 20, Users(expected), nil)
	u, err := p.LookupUser("user")
	assert.NoError(t, err)
	if assert.NotNil(t, u) {
		assert.Equal(t, expected, u)
	}
}

func Test_FakeProvider_LookupUser_DoesNotExist(t *testing.T) {
	p := MakeFakeProvider(501, 20, nil, nil)
	u, err := p.LookupUser("user")
	assert.Error(t, err)
	assert.Nil(t, u)
}

func Test_FakeProvider_LookupUserId_Exists(t *testing.T) {
	expected := &user.User{
		Uid: "502",
	}
	p := MakeFakeProvider(501, 20, Users(expected), nil)
	u, err := p.LookupUserId("502")
	assert.NoError(t, err)
	if assert.NotNil(t, u) {
		assert.Equal(t, expected, u)
	}
}

func Test_FakeProvider_LookupUserId_DoesNotExist(t *testing.T) {
	p := MakeFakeProvider(501, 20, nil, nil)
	u, err := p.LookupUserId("502")
	assert.Error(t, err)
	assert.Nil(t, u)
}

func Test_FakeProvider_LookupGroup_Exists(t *testing.T) {
	expected := &group.Group{
		Name: "group",
	}
	p := MakeFakeProvider(501, 20, nil, Groups(expected))
	g, err := p.LookupGroup("group")
	assert.NoError(t, err)
	if assert.NotNil(t, g) {
		assert.Equal(t, expected, g)
	}
}

func Test_FakeProvider_LookupGroup_DoesNotExist(t *testing.T) {
	p := MakeFakeProvider(501, 20, nil, nil)
	g, err := p.LookupGroup("group")
	assert.Error(t, err)
	assert.Nil(t, g)
}

func Test_FakeProvider_LookupGroupId_Exists(t *testing.T) {
	expected := &group.Group{
		Gid: "502",
	}
	p := MakeFakeProvider(501, 20, nil, Groups(expected))
	g, err := p.LookupGroupId("502")
	assert.NoError(t, err)
	if assert.NotNil(t, g) {
		assert.Equal(t, expected, g)
	}
}

func Test_FakeProvider_LookupGroupId_DoesNotExist(t *testing.T) {
	p := MakeFakeProvider(501, 20, nil, nil)
	g, err := p.LookupGroupId("502")
	assert.Error(t, err)
	assert.Nil(t, g)
}

func Test_FakeProvider_Geteuid(t *testing.T) {
	p := MakeFakeProvider(501, 20, nil, nil)
	assert.Equal(t, 501, p.Geteuid())
}

func Test_FakeProvider_Getegid(t *testing.T) {
	p := MakeFakeProvider(501, 20, nil, nil)
	assert.Equal(t, 20, p.Getegid())
}
