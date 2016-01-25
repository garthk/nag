package userinfo

import (
	"fmt"
	"os/user"

	"github.com/milosgajdos83/servpeek/utils/group"
)

// MakeFakeProvider creates a fake user info provider.
func MakeFakeProvider(euid, egid int, users []*user.User, groups []*group.Group) Provider {
	p := fakeUserInfoProvider{
		euid:         euid,
		egid:         egid,
		usersByName:  make(map[string]*user.User),
		usersById:    make(map[string]*user.User),
		groupsByName: make(map[string]*group.Group),
		groupsById:   make(map[string]*group.Group),
	}

	for _, u := range users {
		p.usersByName[u.Username] = u
		p.usersById[u.Uid] = u
	}

	for _, g := range groups {
		p.groupsByName[g.Name] = g
		p.groupsById[g.Gid] = g
	}

	return p
}

type fakeUserInfoProvider struct {
	usersByName  map[string]*user.User
	usersById    map[string]*user.User
	groupsByName map[string]*group.Group
	groupsById   map[string]*group.Group
	euid         int
	egid         int
}

func (p fakeUserInfoProvider) LookupUser(username string) (*user.User, error) {
	if u := p.usersByName[username]; u == nil {
		return nil, fmt.Errorf("No user matching username: %#v", username)
	} else {
		return u, nil
	}
}

func (p fakeUserInfoProvider) LookupUserId(uid string) (*user.User, error) {
	if u := p.usersById[uid]; u == nil {
		return nil, fmt.Errorf("No user matching uid: %#v", uid)
	} else {
		return u, nil
	}
}

func (p fakeUserInfoProvider) LookupGroup(groupname string) (*group.Group, error) {
	if g := p.groupsByName[groupname]; g == nil {
		return nil, fmt.Errorf("No group matching groupname: %#v", groupname)
	} else {
		return g, nil
	}
}

func (p fakeUserInfoProvider) LookupGroupId(gid string) (*group.Group, error) {
	if g := p.groupsById[gid]; g == nil {
		return nil, fmt.Errorf("No user matching gid: %#v", gid)
	} else {
		return g, nil
	}
}

func (p fakeUserInfoProvider) Geteuid() int {
	return p.euid
}

func (p fakeUserInfoProvider) Getegid() int {
	return p.egid
}

func Users(users ...*user.User) []*user.User {
	return users
}

func Groups(groups ...*group.Group) []*group.Group {
	return groups
}
