package naglib

type PluginContext struct {
	NagiosConfig        *NagiosConfig
	ExitStatusTreatment *ExitStatusTreatment
	Messages            []Message
}

func (pc *PluginContext) Add(severity PluginStatus, message string) {
	pc.Messages = append(pc.Messages, Message{
		Message:  message,
		Severity: severity,
	})
}

type Message struct {
	Message  string
	Severity PluginStatus
	Err      error
}
