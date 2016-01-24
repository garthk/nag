package naglib

import (
	"time"
)

type NagiosConfig struct {
	RunAsUser      string
	RunAsGroup     string
	AllowArguments bool
	CommandPrefix  string
	CommandTimeout time.Duration
	NonFatalErrors []error
	Commands       map[string]string
}

const DEFAULT_COMMAND_TIMEOUT = 60 // nrpe.c

func NewNagiosConfig() *NagiosConfig {
	return &NagiosConfig{
		Commands:       make(map[string]string),
		CommandTimeout: DEFAULT_COMMAND_TIMEOUT * time.Second,
	}
}
