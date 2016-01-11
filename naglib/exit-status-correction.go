package naglib

type ExitStatusTreatment struct {
	CriticalWarnings bool
	CriticalUnknowns bool
	CriticalExcess   bool
	TolerantExcess   bool
}

func ExitStatusToPluginStatus(exitStatus int, treatment *ExitStatusTreatment) PluginStatus {
	switch {
	case exitStatus == 0:
		return OK

	case exitStatus == 1:
		if treatment.CriticalWarnings {
			return CRITICAL
		} else {
			return WARNING
		}

	case exitStatus == 2:
		return CRITICAL

	case exitStatus == 3:
		if treatment.CriticalUnknowns {
			return CRITICAL
		} else {
			return UNKNOWN
		}
	}

	if treatment.CriticalExcess {
		return CRITICAL
	} else {
		return UNKNOWN
	}
}
