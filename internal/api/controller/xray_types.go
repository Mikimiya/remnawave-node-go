package controller

import (
	"github.com/hteppl/remnawave-node-go/internal/utils"
	"github.com/hteppl/remnawave-node-go/internal/xray"
)

const (
	APIPort = 61012
)

var (
	msgRequestAlreadyInProgress  = "Request already in progress"
	msgUnsupportedVersion        = "Unsupported Remnawave version. Please, upgrade Remnawave to version v2.3.x or higher"
	logStartAlreadyInProgress    = "Start request already in progress, rejecting duplicate"
	logFailedToParseStartRequest = "Failed to parse start request"
	logRestartRequired           = "Restart required - proceeding with xray core restart"
	logForceRestartRequested     = "Force restart requested"
	logFailedToExtractUsers      = "Failed to extract users from config"
	logFailedToMarshalConfig     = "Failed to marshal xray config"
	logFailedToStartXray         = "Failed to start xray core"
	logXrayStartedSuccessfully   = "Xray core started successfully"
	logStopRequested             = "Remnawave requested to stop Xray."
	logFailedToStopXray          = "Failed to stop xray core"
	logXrayStoppedSuccessfully   = "Xray core stopped successfully"
)

type StartRequest struct {
	XrayConfig map[string]interface{} `json:"xrayConfig" binding:"required"`
	Internals  xray.Internals         `json:"internals" binding:"required"`
}

type NodeInformation struct {
	Version string `json:"version"`
}

type SystemInfo struct {
	CpuCores    int    `json:"cpuCores"`
	CpuModel    string `json:"cpuModel"`
	MemoryTotal string `json:"memoryTotal"`
}

type StartResponse struct {
	IsStarted       bool            `json:"isStarted"`
	Version         *string         `json:"version"`
	Error           *string         `json:"error"`
	SystemInfo      *SystemInfo     `json:"systemInformation"`
	NodeInformation NodeInformation `json:"nodeInformation"`
}

type StopResponse struct {
	IsStopped bool `json:"isStopped"`
}

type StatusResponse struct {
	IsRunning bool    `json:"isRunning"`
	Version   *string `json:"version"`
}

type HealthcheckResponse struct {
	IsAlive                  bool    `json:"isAlive"`
	XrayInternalStatusCached bool    `json:"xrayInternalStatusCached"`
	XrayVersion              *string `json:"xrayVersion"`
	NodeVersion              string  `json:"nodeVersion"`
}

func getSystemInfo() SystemInfo {
	return SystemInfo{
		CpuCores:    utils.GetCPUCores(),
		CpuModel:    utils.GetCPUModel(),
		MemoryTotal: utils.GetTotalMemory(),
	}
}
