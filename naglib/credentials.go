package naglib

import (
	"errors"
	"fmt"
	"log"
	"os/user"
	"regexp"
	"strconv"
	"syscall"

	"github.com/garthk/nag/pkg/user-info-shim"
	"github.com/milosgajdos83/servpeek/utils/group"
)

func getSysProcAttr(provider userinfo.Provider, context *PluginContext) (*syscall.SysProcAttr, error) {
	if context == nil {
		return nil, errors.New("nil context")
	}

	if provider == nil {
		provider = userinfo.GetDefaultProvider()
	}

	cfg := context.NagiosConfig

	newuid, newgid, err := getNewIdentity(provider, cfg.RunAsUser, cfg.RunAsGroup)
	if err != nil {
		context.Add(WARNING, fmt.Sprintf("user/group warning: %s", err.Error()))
	}

	if newuid == uint32(provider.Geteuid()) {
		log.Printf("no need to setuid: already %d", newuid)
		// TODO: check group; figure out edge case
		return nil, nil
	}

	sysprocattr := &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(newuid),
			Gid: uint32(newgid),
		},
	}

	return sysprocattr, nil
}

var allDigits = regexp.MustCompile("^[0-9]+$")

func getNewIdentity(provider userinfo.Provider, userNameOrId, groupNameOrId string) (uid, gid uint32, err error) {
	// Get the new identity.
	// Return uid=euid, gid=egid if we can't or shouldn't change identity.
	// Return uid=N, gid=0 if we should change identity. Note gid always 0, and should be ignored.
	// Return warning!=nil to say why.

	uid = uint32(provider.Geteuid())
	gid = uint32(provider.Getegid())
	err = nil

	if userNameOrId == "" {
		userNameOrId = fmt.Sprintf("%d", uid)
	}

	u, err := lookUpUser(provider, userNameOrId)
	if err != nil {
		return
	}

	if groupNameOrId == "" {
		groupNameOrId = u.Gid
	}

	g, err := lookupGroup(provider, groupNameOrId)
	if err != nil {
		return
	}

	puid, err := strconv.ParseUint(u.Uid, 10, 32)
	if err != nil {
		err = fmt.Errorf("bad retrieved Uid: %s", err.Error())
		return
	}

	pgid, err := strconv.ParseUint(g.Gid, 10, 32)
	if err != nil {
		err = fmt.Errorf("bad retrieved Gid: %s", err.Error())
		return
	}

	uid = uint32(puid)
	gid = uint32(pgid)

	return
}

type userLookup func(string) (*user.User, error)

func lookUpUser(provider userinfo.Provider, userNameOrId string) (u *user.User, err error) {
	var fn userLookup
	var what string

	if allDigits.MatchString(userNameOrId) {
		fn = provider.LookupUserId
		what = "user ID"
	} else {
		fn = provider.LookupUser
		what = "user"
	}

	u, err = fn(userNameOrId)
	if err != nil {
		err = errors.New(fmt.Sprintf("can't look up %s %#v: %s", what, userNameOrId, err.Error()))
	}

	return
}

type groupLookup func(string) (*group.Group, error)

func lookupGroup(provider userinfo.Provider, groupNameOrId string) (g *group.Group, err error) {
	var fn groupLookup
	var what string

	if allDigits.MatchString(groupNameOrId) {
		fn = provider.LookupGroupId
		what = "group ID"
	} else {
		fn = provider.LookupGroup
		what = "group"
	}

	g, err = fn(groupNameOrId)
	if err != nil {
		err = errors.New(fmt.Sprintf("can't look up %s %#v: %s", what, groupNameOrId, err.Error()))
	}

	return
}
