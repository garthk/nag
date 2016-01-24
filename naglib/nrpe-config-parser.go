package naglib

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/garthk/nag/pkg/readable-fs"
)

// NRPE configuration parser.
//
// More compatible with NRPE than you'd really want, but not 100%:
//
// * I couldn't stomach matching its misbehaviour on lines >MAX_INPUT_BUFFER
// * It ignores configuration values not known to it
//
// For reference, please see:
//
// http://sourceforge.net/p/nagios/nrpe/ci/master/tree/src/nrpe.c
// http://sourceforge.net/p/nagios/nrpe/ci/master/tree/include/common.h
// http://sourceforge.net/p/nagios/nrpe/ci/master/tree/sample-config/nrpe.cfg.in

const MAX_INPUT_BUFFER = 2048      // common.h
const MAX_FILENAME_LENGTH = 256    // common.h

// WARNING: you MUST maintain these in sync with ../pkg/fake-readable-fs/types.go


const NRPE_NO_VAR_ERR = "No variable name specified in config file '%s' - Line '%d'"
const NRPE_NO_VAL_ERR = "No variable value specified in config file '%s' - Line '%d'"
const NRPE_DIR_OPEN_ERR = "Could not open config directory '%s' for reading."
const NRPE_FILE_OPEN_ERR = "Unable to open config file '%s' for reading"
const NRPE_BAD_TIMEOUT = "Invalid command_timeout specified in config file '%s' - Line %d"
const NRPE_BAD_COMMAND = "Invalid command specified in config file '%s' - Line %d"

func ParseConfig(filename string) (*NagiosConfig, error) {
	return parseConfig(nil, rfs.Reality(), filename)
}

func parseConfig(cfg *NagiosConfig, fs rfs.ReadableFileSystem, filename string) (*NagiosConfig, error) {
	if cfg == nil {
		cfg = NewNagiosConfig()
	}

	rawbytes, err := fs.ReadFile(filename)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(NRPE_FILE_OPEN_ERR, filename))
	}

	buffer := bytes.NewBuffer(rawbytes)
	scanner := bufio.NewScanner(buffer)
	lineno := 0

	for scanner.Scan() {
		// I'm deliberately using the methodology from read_config_file in
		// nrpe.c to ensure compatibility, and many variable names from it
		// to make it easy to verify compatibility.

		lineno = lineno + 1
		input_line := scanner.Text()
		input_line = strings.TrimSpace(input_line)
		eqloc := strings.Index(input_line, "=")

		switch {
		case len(input_line) == 0: // empty line
			continue
		case input_line[0] == '#': // comment
			continue
		case eqloc < 0: // no equals sign
			return errorResult(NRPE_NO_VAL_ERR, filename, lineno)
		case len(input_line) == 1: // =
			return errorResult(NRPE_NO_VAR_ERR, filename, lineno)
		case eqloc == 0: // =value
			return errorResult(NRPE_NO_VAL_ERR, filename, lineno)
		case eqloc == len(input_line)-1:
			return errorResult(NRPE_NO_VAL_ERR, filename, lineno)
		}

		varname := input_line[:eqloc]
		varvalue := input_line[eqloc+1:]

		switch varname {
		case "nrpe_user":
			cfg.RunAsUser = varvalue

		case "nrpe_group":
			cfg.RunAsGroup = varvalue

		// nrpe.c: allow_arguments=atoi(varvalue)==1?TRUE:FALSE
		// ... so 1 and 1Literally are TRUE
		// ... and 0, 0hope, and 16CandlesThereOnMyWall are FALSE

		case "dont_blame_nrpe":
			if badAtoi(varvalue) == 1 {
				cfg.AllowArguments = true
			} else {
				cfg.AllowArguments = false
			}

		case "command_prefix":
			cfg.CommandPrefix = varvalue

		// nrpe.c: command_timeout=atoi(varvalue)
		// ... so 500Miles and 100000Morriseys are OK

		case "command_timeout":
			command_timeout := badAtoi(varvalue)

			if command_timeout < 1 {
				return errorResult(NRPE_BAD_TIMEOUT, filename, lineno)
			}

			cfg.CommandTimeout = time.Duration(command_timeout) * time.Second

		// nrpe.c calls chdir("/") when run, unless it's from inetd, and
		// doesn't treat include or include_dir arguments as relative to
		// the file being processed. I can't see how relative arguments
		// can be reliable.

		case "include":
			_, err := parseConfig(cfg, fs, varvalue)
			if err != nil {
				cfg.NonFatalErrors = append(cfg.NonFatalErrors, err)
			}

		case "include_file":
			_, err := parseConfig(cfg, fs, varvalue)
			if err != nil {
				cfg.NonFatalErrors = append(cfg.NonFatalErrors, err)
			}

		case "include_dir":
			_, err := parseConfigDir(cfg, fs, varvalue)
			if err != nil {
				cfg.NonFatalErrors = append(cfg.NonFatalErrors, err)
			}

		default:
			if len(varname) >= 8 && varname[:8] == "command[" {
				command := varname[8:]
				end := strings.Index(command, "]")
				if end < 1 {
					return errorResult(NRPE_BAD_COMMAND, filename, lineno)
				}
				command = command[:end]
				cfg.Commands[command] = varvalue
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func errorResult(format string, filename string, lineno int) (*NagiosConfig, error) {
	return nil, errors.New(fmt.Sprintf(format, filename, lineno))
}

func badAtoi(str string) int {
	digits := leadingDigits(str)

	if len(digits) == 0 {
		return 0
	}

	i, err := strconv.Atoi(digits)
	if err != nil {
		panic(err)
	}

	return i
}

func leadingDigits(str string) string {
	r, _ := regexp.Compile("^[0-9]+")
	return r.FindString(str)
}

func parseConfigDir(cfg *NagiosConfig, fs rfs.ReadableFileSystem, dirname string) (*NagiosConfig, error) {
	files, err := fs.ReadDir(dirname)

	if err != nil {
		return nil, errors.New(fmt.Sprintf(NRPE_DIR_OPEN_ERR, dirname))
	}

	for i := range files {
		basename := files[i].Name()
		filename := filepath.Join(dirname, basename)

		if filepath.Ext(filename) != ".cfg" {
			continue
		}

		_, err := parseConfig(cfg, fs, filepath.Join(dirname, files[i].Name()))
		if err != nil {
			cfg.NonFatalErrors = append(cfg.NonFatalErrors, err)
		}
	}

	return cfg, nil
}
