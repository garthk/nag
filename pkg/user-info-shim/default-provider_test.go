package userinfo

import (
	"os"
	"os/user"
	"testing"

	"github.com/milosgajdos83/servpeek/utils/group"
	"github.com/stretchr/testify/assert"
)

func Test_DefaultProvider_LookupUser(t *testing.T) {
	current := getCurrentUser(t)
	uip := GetDefaultProvider()

	u1, err := user.Lookup(current.Username)
	if err != nil {
		t.Fatalf("Can't look up current user directly: %s", err.Error())
	}

	u2, err := uip.LookupUser(current.Username)
	if err != nil {
		t.Fatalf("Can't look up current user via provider: %s", err.Error())
	}

	assert.Equal(t, u1.Username, u2.Username)
}

func Test_DefaultProvider_LookupUserId(t *testing.T) {
	current := getCurrentUser(t)
	uip := GetDefaultProvider()

	u1, err := user.LookupId(current.Uid)
	if err != nil {
		t.Fatalf("Can't look up current user directly: %s", err.Error())
	}

	u2, err := uip.LookupUserId(current.Uid)
	if err != nil {
		t.Fatalf("Can't look up current user via provider: %s", err.Error())
	}

	assert.Equal(t, u1.Uid, u2.Uid)
}

func Test_DefaultProvider_LookupGroupId(t *testing.T) {
	current := getCurrentUser(t)
	uip := GetDefaultProvider()

	g1, err := group.LookupID(current.Gid)
	if err != nil {
		t.Fatalf("Can't look up current group directly: %s", err.Error())
	}

	g2, err := uip.LookupGroupId(current.Gid)
	if err != nil {
		t.Fatalf("Can't look up current group via provider: %s", err.Error())
	}

	assert.Equal(t, g1.Gid, g2.Gid)
}

func Test_DefaultProvider_LookupGroup(t *testing.T) {
	current := getCurrentUser(t)
	uip := GetDefaultProvider()

	g1, err := group.LookupID(current.Gid)
	if err != nil {
		t.Fatalf("Can't look up current group by ID directly: %s", err.Error())
	}

	g1b, err := group.Lookup(g1.Name)
	if err != nil {
		t.Fatalf("Can't look up current group by name directly: %s", err.Error())
	}

	g2, err := uip.LookupGroup(g1.Name)
	if err != nil {
		t.Fatalf("Can't look up current user via provider: %s", err.Error())
	}

	assert.Equal(t, g1.Name, g1b.Name, "group module failure")
	assert.Equal(t, g1b.Name, g2.Name, "shim failure")
}

func Test_DefaultProvider_Geteuid(t *testing.T) {
	uip := GetDefaultProvider()
	assert.Equal(t, os.Geteuid(), uip.Geteuid())
}

func Test_DefaultProvider_Getegid(t *testing.T) {
	uip := GetDefaultProvider()
	assert.Equal(t, os.Getegid(), uip.Getegid())
}

func getCurrentUser(t *testing.T) *user.User {
	current, err := user.Current()
	if err != nil {
		t.Fatalf("Can't determine current user for test: %s", err.Error())
	}
	return current
}
