package naglib

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/garthk/nag/pkg/readable-fs"
	"github.com/stretchr/testify/assert"
)

const defaultConfig = `# Interesting parts of default config on Ubuntu:
nrpe_user=nagios
nrpe_group=nagios
dont_blame_nrpe=0
# command_prefix=/usr/bin/sudo
command_timeout=60
include=/etc/nagios/nrpe_local.cfg
include_dir=/etc/nagios/nrpe.d/
`

func Test_ParseConfig_Survival_REALITY(t *testing.T) {
	cfg, err := ParseConfig("\tSlartibartfast") // could _possibly_ fail

	assert.Nil(t, cfg)
	assert.EqualError(t, err, "Unable to open config file '\tSlartibartfast' for reading")
}

func Test_ParseConfig_Survival_DefaultConfig(t *testing.T) {
	fs := rfs.NewFake()
	filename := fixpath("etc/nagios/nrpe.cfg")
	fs.AddFile(filename, rfs.FromString(defaultConfig), 0644)

	cfg, err := testWith(defaultConfig)

	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, "nagios", cfg.RunAsUser)
	assert.Equal(t, "nagios", cfg.RunAsGroup)
	assert.Equal(t, false, cfg.AllowArguments)
	assert.Equal(t, "", cfg.CommandPrefix)
	assert.Equal(t, DEFAULT_COMMAND_TIMEOUT*time.Second, cfg.CommandTimeout)
	assert.Equal(t, 2, len(cfg.NonFatalErrors))
	assert.Equal(t, "Unable to open config file '/etc/nagios/nrpe_local.cfg' for reading", cfg.NonFatalErrors[0].Error())
	assert.Equal(t, "Could not open config directory '/etc/nagios/nrpe.d/' for reading.", cfg.NonFatalErrors[1].Error())
}

var cfgpath string // see init()

func Test_ParseConfig_EmptyLines(t *testing.T) {
	cfg, err := testWith("\n\n\n")

	assert.NotNil(t, cfg)
	assert.Nil(t, err)

	assert.Equal(t, "", cfg.RunAsUser)
	assert.Equal(t, "", cfg.RunAsGroup)
	assert.Equal(t, false, cfg.AllowArguments)
	assert.Equal(t, "", cfg.CommandPrefix)
	assert.Equal(t, DEFAULT_COMMAND_TIMEOUT*time.Second, cfg.CommandTimeout)
	assert.Equal(t, 0, len(cfg.NonFatalErrors))
}

func Test_ParseConfig_CommentLines_OK(t *testing.T) {
	cfg, err := testWith("# this is a comment")

	assert.NotNil(t, cfg)
	assert.Nil(t, err)
}

func Test_ParseConfig_NoValue_Fails(t *testing.T) {
	cfg, err := testWith("foo=")

	assert.Nil(t, cfg)
	assert.EqualError(t, err, fmt.Sprintf(NRPE_NO_VAL_ERR, cfgpath, 1))
}

func Test_ParseConfig_NoEquals_Fails(t *testing.T) {
	cfg, err := testWith("foo")

	assert.Nil(t, cfg)
	assert.EqualError(t, err, fmt.Sprintf(NRPE_NO_VAL_ERR, cfgpath, 1))
}

func Test_ParseConfig_SoloEquals_Fails(t *testing.T) {
	cfg, err := testWith("=")

	assert.Nil(t, cfg)
	assert.EqualError(t, err, fmt.Sprintf(NRPE_NO_VAR_ERR, cfgpath, 1))
}

func Test_ParseConfig_NoVariable_Fails(t *testing.T) {
	cfg, err := testWith("=foo")

	assert.Nil(t, cfg)
	assert.EqualError(t, err, fmt.Sprintf(NRPE_NO_VAL_ERR, cfgpath, 1))
}

func Test_ParseConfig_CommandPrefix(t *testing.T) {
	cfg, err := testWith("command_prefix=sudo")

	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, "sudo", cfg.CommandPrefix)
}

func Test_ParseConfig_MultipleAssignmentOK_LastWins(t *testing.T) {
	cfg, err := testWith(`command_prefix=sudo
	                      command_prefix=nope`)

	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, "nope", cfg.CommandPrefix)
}

func Test_ParseConfig_CommandTimeout_0(t *testing.T) {
	cfg, err := testWith("command_timeout=0")

	assert.Nil(t, cfg)
	assert.EqualError(t, err, fmt.Sprintf(NRPE_BAD_TIMEOUT, cfgpath, 1))
}

func Test_ParseConfig_CommandTimeout_30(t *testing.T) {
	cfg, err := testWith("command_timeout=30")

	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, 30*time.Second, cfg.CommandTimeout)
}

func Test_ParseConfig_CommandTimeout_30hours_Wut(t *testing.T) {
	cfg, err := testWith("command_timeout=30hours")

	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, 30*time.Second, cfg.CommandTimeout)
}

func Test_ParseConfig_NrpeUser(t *testing.T) {
	cfg, err := testWith("nrpe_user=nagios")

	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, "nagios", cfg.RunAsUser)
}

func Test_ParseConfig_NrpeGroup(t *testing.T) {
	cfg, err := testWith("nrpe_group=nagios")

	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, "nagios", cfg.RunAsGroup)
}

func Test_ParseConfig_DontBlameNrpe_0(t *testing.T) {
	cfg, err := testWith("dont_blame_nrpe=0")

	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, false, cfg.AllowArguments)
}

func Test_ParseConfig_DontBlameNrpe_1(t *testing.T) {
	cfg, err := testWith("dont_blame_nrpe=1")

	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, true, cfg.AllowArguments)
}

func Test_ParseConfig_DontBlameNrpe_2_Wut(t *testing.T) {
	cfg, err := testWith("dont_blame_nrpe=2")

	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, false, cfg.AllowArguments)
}

func Test_ParseConfig_DontBlameNrpe_1Food_Wut(t *testing.T) {
	cfg, err := testWith("dont_blame_nrpe=1JobYouHad")

	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, true, cfg.AllowArguments)
}

func Test_ParseConfig_Include(t *testing.T) {
	cfgpath2 := fixpath("/etc/nagios/nrpe_local.cfg")
	cfgpath3 := fixpath("/etc/nagios/nrpe_local2.cfg")

	fs := rfs.NewFake()
	fs.AddFile(cfgpath, rfs.FromString(`include=/etc/nagios/nrpe_local.cfg
	                                    include_file=/etc/nagios/nrpe_local2.cfg`), 0644)
	fs.AddFile(cfgpath2, rfs.FromString("command_timeout=30"), 0644)
	fs.AddFile(cfgpath3, rfs.FromString("command_prefix=sudo"), 0644)

	cfg, err := parseConfig(nil, fs, cfgpath)

	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, 30*time.Second, cfg.CommandTimeout) // via include
	assert.Equal(t, "sudo", cfg.CommandPrefix)          // via include_file
}

func Test_ParseConfig_Include_ErrorsBecomeNonFatal(t *testing.T) {
	cfgpath2 := fixpath("/etc/nagios/nrpe_local.cfg")
	cfgpath3 := fixpath("/etc/nagios/nrpe_local2.cfg")

	fs := rfs.NewFake()
	fs.AddFile(cfgpath, rfs.FromString(`include=/etc/nagios/nrpe_local.cfg
	                                    include_file=/etc/nagios/nrpe_local2.cfg`), 0644)
	fs.AddFile(cfgpath2, rfs.FromString("command_timeout=0"), 0644)
	fs.AddFile(cfgpath3, rfs.FromString("command_prefix="), 0644)

	cfg, err := parseConfig(nil, fs, cfgpath)

	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, 60*time.Second, cfg.CommandTimeout) // default still there
	assert.Equal(t, 2, len(cfg.NonFatalErrors))
	assert.Equal(t, fmt.Sprintf(NRPE_BAD_TIMEOUT, cfgpath2, 1), cfg.NonFatalErrors[0].Error())
	assert.Equal(t, fmt.Sprintf(NRPE_NO_VAL_ERR, cfgpath3, 1), cfg.NonFatalErrors[1].Error())
}

func Test_ParseConfig_IncludeDir(t *testing.T) {
	cfgpath2 := fixpath("/etc/nagios/nrpe.d/extra.cfg")

	fs := rfs.NewFake()
	fs.AddFile(cfgpath, rfs.FromString("include_dir=/etc/nagios/nrpe.d/"), 0644)
	fs.AddFile(cfgpath2, rfs.FromString("command_timeout=30"), 0644)

	cfg, err := parseConfig(nil, fs, cfgpath)

	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, 30*time.Second, cfg.CommandTimeout)
	assert.Equal(t, 0, len(cfg.NonFatalErrors))
}

func Test_ParseConfig_IncludeDir_SkipNonCfg(t *testing.T) {
	cfgpath2 := fixpath("/etc/nagios/nrpe.d/extra.not-cfg")

	fs := rfs.NewFake()
	fs.AddFile(cfgpath, rfs.FromString("include_dir=/etc/nagios/nrpe.d/"), 0644)
	fs.AddFile(cfgpath2, rfs.FromString("command_timeout=0"), 0644)

	cfg, err := parseConfig(nil, fs, cfgpath)

	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(cfg.NonFatalErrors)) // 1 if cfgpath2 read
}

func Test_ParseConfig_IncludeDir_NoStopOnErrors(t *testing.T) {
	cfgpath2 := fixpath("/etc/nagios/nrpe.d/02.cfg") // who said they were sorted?
	cfgpath3 := fixpath("/etc/nagios/nrpe.d/01.cfg")

	fs := rfs.NewFake()
	fs.AddFile(cfgpath, rfs.FromString("include_dir=/etc/nagios/nrpe.d/"), 0644)
	fs.AddFile(cfgpath2, rfs.FromString("command_timeout=0"), 0644)
	fs.AddFile(cfgpath3, rfs.FromString("command_timeout=30"), 0644)

	cfg, err := parseConfig(nil, fs, cfgpath)

	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, 30*time.Second, cfg.CommandTimeout)
	assert.Equal(t, 1, len(cfg.NonFatalErrors))
	assert.Equal(t, fmt.Sprintf(NRPE_BAD_TIMEOUT, cfgpath2, 1), cfg.NonFatalErrors[0].Error())
}

func Test_ParseConfig_Command(t *testing.T) {
	cfg, err := testWith("command[ok]=/bin/true")

	assert.NotNil(t, cfg)
	assert.Nil(t, err)
	assert.Equal(t, "/bin/true", cfg.Commands["ok"])
}

func Test_ParseConfig_No_Command(t *testing.T) {
	cfg, err := testWith("command[]=/bin/true")

	assert.Nil(t, cfg)
	assert.EqualError(t, err, fmt.Sprintf(NRPE_BAD_COMMAND, cfgpath, 1))
}

// TODO: warnings on lines longer than MAX_INPUT_BUFFER
// TODO: warnings on joined filenames longer than MAX_FILENAME_LENGTH
// TODO: warnings on relative arguments to include, include_dir
// TODO: warnings on atoi madness
// TODO: warnings on dont_blame_nrpe=2 wut

func Test_badAtoi_nodigits(t *testing.T) {
	i := badAtoi("")
	assert.Equal(t, 0, i)
}

func init() {
	cfgpath = fixpath("/etc/nagios/nrpe.cfg")
}

func testWith(contents string) (*NagiosConfig, error) {
	fs := rfs.NewFake()
	fs.AddFile(cfgpath, rfs.FromString(contents), 0644)
	return parseConfig(nil, fs, cfgpath)
}

func fixpath(path string) string {
	root, err := rfs.RootPath()
	if err != nil {
		panic(err)
	}

	inparts := strings.Split(path, "/")
	if len(inparts) > 0 && inparts[0] == "" {
		inparts = inparts[1:]
	}

	parts := make([]string, 0)
	parts = append(parts, root)
	parts = append(parts, inparts...)

	return filepath.Join(parts...)
}
