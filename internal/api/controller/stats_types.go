package controller

var (
	logFailedToParseInboundStatsRequest  = "Failed to parse get-inbound-stats request"
	logFailedToParseOutboundStatsRequest = "Failed to parse get-outbound-stats request"
)

type ResetRequest struct {
	Reset bool `json:"reset"`
}

type UsernameRequest struct {
	Username string `json:"username" binding:"required"`
}

type UserIDRequest struct {
	UserID string `json:"userId" binding:"required"`
}

type TagResetRequest struct {
	Tag   string `json:"tag" binding:"required"`
	Reset bool   `json:"reset"`
}

type SystemStatsResponse struct {
	NumGoroutine int    `json:"numGoroutine"`
	NumGC        uint32 `json:"numGC"`
	Alloc        uint64 `json:"alloc"`
	TotalAlloc   uint64 `json:"totalAlloc"`
	Sys          uint64 `json:"sys"`
	Mallocs      uint64 `json:"mallocs"`
	Frees        uint64 `json:"frees"`
	LiveObjects  uint64 `json:"liveObjects"`
	PauseTotalNs uint64 `json:"pauseTotalNs"`
	Uptime       int64  `json:"uptime"`
}

type UserStats struct {
	Username string `json:"username"`
	Uplink   int64  `json:"uplink"`
	Downlink int64  `json:"downlink"`
}

type UsersStatsResponse struct {
	Users []UserStats `json:"users"`
}

type UserOnlineResponse struct {
	IsOnline bool `json:"isOnline"`
}

type UserIPListResponse struct {
	IPs []string `json:"ips"`
}

type InboundStatsResponse struct {
	Inbound  string `json:"inbound"`
	Uplink   int64  `json:"uplink"`
	Downlink int64  `json:"downlink"`
}

type OutboundStatsResponse struct {
	Outbound string `json:"outbound"`
	Uplink   int64  `json:"uplink"`
	Downlink int64  `json:"downlink"`
}

type InboundEntry struct {
	Inbound  string `json:"inbound"`
	Uplink   int64  `json:"uplink"`
	Downlink int64  `json:"downlink"`
}

type AllInboundsStatsResponse struct {
	Inbounds []InboundEntry `json:"inbounds"`
}

type OutboundEntry struct {
	Outbound string `json:"outbound"`
	Uplink   int64  `json:"uplink"`
	Downlink int64  `json:"downlink"`
}

type AllOutboundsStatsResponse struct {
	Outbounds []OutboundEntry `json:"outbounds"`
}

type CombinedStatsResponse struct {
	Inbounds  []InboundEntry  `json:"inbounds"`
	Outbounds []OutboundEntry `json:"outbounds"`
}
