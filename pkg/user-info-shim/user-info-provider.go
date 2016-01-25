// Package userinfo provides a shim to, and mock of, user information
// methods in the os and os/user packages.
package userinfo // import "github.com/garthk/nag/pkg/user-info-shim"

import (
	"os/user"

	"github.com/milosgajdos83/servpeek/utils/group"
)

type Provider interface {
	Getegid() int                                      // os.Getegid
	Geteuid() int                                      // os.Geteuid
	LookupUser(username string) (*user.User, error)    // user.Lookup
	LookupUserId(uid string) (*user.User, error)       // user.LookupId
	LookupGroup(username string) (*group.Group, error) // group.Lookup
	LookupGroupId(uid string) (*group.Group, error)    // group.LookupId
}
