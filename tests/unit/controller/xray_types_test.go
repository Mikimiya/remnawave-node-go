package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hteppl/remnawave-node-go/internal/api/controller"
	"github.com/hteppl/remnawave-node-go/internal/logger"
	"github.com/hteppl/remnawave-node-go/internal/xray"
)

func TestGetSystemInfo(t *testing.T) {
	log := logger.New(logger.Config{Level: logger.LevelError, Format: logger.FormatJSON})
	core := xray.NewCore(log)
	configMgr := xray.NewConfigManager(log)

	xrayCtrl := controller.NewXrayController(core, configMgr, log)

	router := gin.New()
	group := router.Group("/node/xray")
	xrayCtrl.RegisterRoutes(group)

	// Start xray to get a response that includes SystemInfo
	startReq := map[string]interface{}{
		"xrayConfig": map[string]interface{}{
			"log": map[string]interface{}{"loglevel": "warning"},
			"inbounds": []interface{}{
				map[string]interface{}{
					"tag": "vless-in", "port": 10000, "protocol": "vless",
					"settings":       map[string]interface{}{"clients": []interface{}{}, "decryption": "none"},
					"streamSettings": map[string]interface{}{"network": "tcp"},
				},
			},
			"outbounds": []interface{}{
				map[string]interface{}{"tag": "direct", "protocol": "freedom"},
			},
		},
		"internals": map[string]interface{}{
			"forceRestart": false,
			"hashes": map[string]interface{}{
				"emptyConfig": "abc123",
				"inbounds":    []interface{}{},
			},
		},
	}

	req := httptest.NewRequest("POST", "/node/xray/start", jsonBody(t, startReq))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Response struct {
			IsStarted  bool `json:"isStarted"`
			SystemInfo *struct {
				CpuCores    int    `json:"cpuCores"`
				CpuModel    string `json:"cpuModel"`
				MemoryTotal string `json:"memoryTotal"`
			} `json:"systemInformation"`
		} `json:"response"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response.Response.IsStarted)
	require.NotNil(t, response.Response.SystemInfo)
	assert.Greater(t, response.Response.SystemInfo.CpuCores, 0, "CpuCores should be non-zero")
	assert.NotEmpty(t, response.Response.SystemInfo.CpuModel, "CpuModel should not be empty")
	assert.NotEmpty(t, response.Response.SystemInfo.MemoryTotal, "MemoryTotal should not be empty")
}
