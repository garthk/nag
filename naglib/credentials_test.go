package naglib

import (
	"os/user"
	"testing"

	userinfo "github.com/garthk/nag/pkg/user-info-shim"
	"github.com/milosgajdos83/servpeek/utils/group"
	"github.com/stretchr/testify/assert"
)

var user501 = &user.User{
	Username: "user",
	Uid:      "501",
	Gid:      "20",
}

var user502 = &user.User{
	Username: "user2",
	Uid:      "502",
	Gid:      "12",
}

var staff20 = &group.Group{
	Name: "staff",
	Gid:  "20",
}

var everyone12 = &group.Group{
	Name: "everyone",
	Gid:  "12",
}

var wheel0 = &group.Group{
	Name: "wheel",
	Gid:  "0",
}

var users = []*user.User{user501, user502}
var groups = []*group.Group{staff20, everyone12, wheel0}

func Test_getNewIdentity_NoUser_NoGroup(t *testing.T) {
	p := userinfo.MakeFakeProvider(501, 20, users, groups)

	uid, gid, warning := getNewIdentity(p, "", "")

	assert.Equal(t, uint32(501), uid)
	assert.Equal(t, uint32(20), gid)
	assert.NoError(t, warning)
}

func Test_getNewIdentity_CurrentUser_NoGroup(t *testing.T) {
	p := userinfo.MakeFakeProvider(501, 20, users, groups)

	uid, gid, warning := getNewIdentity(p, "user", "")

	assert.Equal(t, uint32(501), uid)
	assert.Equal(t, uint32(20), gid)
	assert.NoError(t, warning)
}

func Test_getNewIdentity_NoUser_CurrentGroup(t *testing.T) {
	p := userinfo.MakeFakeProvider(501, 20, users, groups)

	uid, gid, warning := getNewIdentity(p, "", "staff")

	assert.Equal(t, uint32(501), uid)
	assert.Equal(t, uint32(20), gid)
	assert.NoError(t, warning)
}

func Test_getNewIdentity_DifferentUserId_ImplicitGroup(t *testing.T) {
	p := userinfo.MakeFakeProvider(501, 20, users, groups)

	uid, gid, warning := getNewIdentity(p, "502", "")

	assert.Equal(t, uint32(502), uid)
	assert.Equal(t, uint32(12), gid)
	assert.NoError(t, warning)
}

func Test_getNewIdentity_DifferentUserId_WithExplicitGroup(t *testing.T) {
	p := userinfo.MakeFakeProvider(501, 20, users, groups)

	uid, gid, warning := getNewIdentity(p, "502", "20")

	assert.Equal(t, uint32(502), uid)
	assert.Equal(t, uint32(20), gid)
	assert.NoError(t, warning)
}

func Test_getNewIdentity_DifferentGroup(t *testing.T) {
	p := userinfo.MakeFakeProvider(501, 20, users, groups)

	uid, gid, warning := getNewIdentity(p, "", "wheel")

	assert.Equal(t, uint32(501), uid)
	assert.Equal(t, uint32(0), gid)
	assert.NoError(t, warning)
}

func Test_getNewIdentity_UnknownUserId(t *testing.T) {
	p := userinfo.MakeFakeProvider(501, 20, users, groups)

	uid, gid, warning := getNewIdentity(p, "503", "")

	assert.Equal(t, uint32(501), uid)
	assert.Equal(t, uint32(20), gid)
	if assert.Error(t, warning) {
		assert.Equal(t, `can't look up user ID "503": No user matching uid: "503"`, warning.Error())
	}
}

func Test_getNewIdentity_UnknownUserName(t *testing.T) {
	p := userinfo.MakeFakeProvider(501, 20, users, groups)

	uid, gid, warning := getNewIdentity(p, "banana", "")

	assert.Equal(t, uint32(501), uid)
	assert.Equal(t, uint32(20), gid)
	if assert.Error(t, warning) {
		assert.Equal(t, `can't look up user "banana": No user matching username: "banana"`, warning.Error())
	}
}

// This should never happen, but with one API using strings and another
// using numbers I feel like I should check it out, anyway.
func Test_getNewIdentity_BadKnownUserId(t *testing.T) {
	u := &user.User{
		Username: "user",
		Uid:      "4294967296",
		Gid:      "20",
	}
	p := userinfo.MakeFakeProvider(501, 20, userinfo.Users(user501, u), groups)

	uid, gid, warning := getNewIdentity(p, "4294967296", "")

	assert.Equal(t, uint32(501), uid)
	assert.Equal(t, uint32(20), gid)
	if assert.Error(t, warning) {
		assert.Equal(t, `bad retrieved Uid: strconv.ParseUint: parsing "4294967296": value out of range`, warning.Error())
	}
}

// This should never happen, but with one API using strings and another
// using numbers I feel like I should check it out, anyway.
func Test_getNewIdentity_BadKnownGroupId(t *testing.T) {
	u := &user.User{
		Username: "user",
		Uid:      "501",
		Gid:      "4294967296",
	}
	g := &group.Group{
		Name: "bad",
		Gid:  "4294967296",
	}
	p := userinfo.MakeFakeProvider(501, 20, userinfo.Users(u), userinfo.Groups(g))

	uid, gid, warning := getNewIdentity(p, "501", "4294967296")

	assert.Equal(t, uint32(501), uid)
	assert.Equal(t, uint32(20), gid)
	if assert.Error(t, warning) {
		assert.Equal(t, `bad retrieved Gid: strconv.ParseUint: parsing "4294967296": value out of range`, warning.Error())
	}
}

func Test_getNewIdentity_UnknownUsername(t *testing.T) {
	p := userinfo.MakeFakeProvider(501, 20, users, groups)

	uid, gid, warning := getNewIdentity(p, "nagios", "")

	assert.Equal(t, uint32(501), uid)
	assert.Equal(t, uint32(20), gid)
	if assert.Error(t, warning) {
		assert.Equal(t, `can't look up user "nagios": No user matching username: "nagios"`, warning.Error())
	}
}

func Test_getNewIdentity_UnknownGroupname(t *testing.T) {
	p := userinfo.MakeFakeProvider(501, 20, users, groups)

	uid, gid, warning := getNewIdentity(p, "", "nagios")

	assert.Equal(t, uint32(501), uid)
	assert.Equal(t, uint32(20), gid)
	if assert.Error(t, warning) {
		assert.Equal(t, `can't look up group "nagios": No group matching groupname: "nagios"`, warning.Error())
	}
}
