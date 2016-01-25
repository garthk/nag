package userinfo

import (
	"os"
	"os/user"

	"github.com/milosgajdos83/servpeek/utils/group"
)

// GetDefaultProvider returns a shim to functions from the Go standard library.
func GetDefaultProvider() Provider {
	return defaultProvider{}
}

type defaultProvider struct{}

func (p defaultProvider) LookupUser(username string) (*user.User, error) {
	return user.Lookup(username)
}

func (p defaultProvider) LookupUserId(uid string) (*user.User, error) {
	return user.LookupId(uid)
}

func (p defaultProvider) LookupGroup(groupname string) (*group.Group, error) {
	return group.Lookup(groupname)
}

func (p defaultProvider) LookupGroupId(gid string) (*group.Group, error) {
	return group.LookupID(gid)
}

func (p defaultProvider) Geteuid() int {
	return os.Geteuid()
}

func (p defaultProvider) Getegid() int {
	return os.Getegid()
}
