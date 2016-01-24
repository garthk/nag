package naglib

type PluginResult struct {
	OriginalOutput string
	Output         string
	Status         PluginStatus
}
