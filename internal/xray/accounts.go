package xray

import (
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/proxy/shadowsocks"
	"github.com/xtls/xray-core/proxy/trojan"
	"github.com/xtls/xray-core/proxy/vless"
)

type CipherType int32

const (
	CipherTypeUnknown           CipherType = 0
	CipherTypeAES128GCM         CipherType = 5
	CipherTypeAES256GCM         CipherType = 6
	CipherTypeCHACHA20POLY1305  CipherType = 7
	CipherTypeXCHACHA20POLY1305 CipherType = 8
	CipherTypeNone              CipherType = 9
)

func BuildVlessUser(email, uuid, flow string, level uint32) *protocol.User {
	vlessAccount := &vless.Account{
		Id:   uuid,
		Flow: flow,
	}

	return &protocol.User{
		Level:   level,
		Email:   email,
		Account: serial.ToTypedMessage(vlessAccount),
	}
}

func BuildTrojanUser(email, password string, level uint32) *protocol.User {
	trojanAccount := &trojan.Account{
		Password: password,
	}

	return &protocol.User{
		Level:   level,
		Email:   email,
		Account: serial.ToTypedMessage(trojanAccount),
	}
}

func BuildShadowsocksUser(email, password string, cipherType CipherType, ivCheck bool, level uint32) *protocol.User {
	ssAccount := &shadowsocks.Account{
		Password:   password,
		CipherType: shadowsocks.CipherType(cipherType),
		IvCheck:    ivCheck,
	}

	return &protocol.User{
		Level:   level,
		Email:   email,
		Account: serial.ToTypedMessage(ssAccount),
	}
}

type UserData struct {
	UserID         string
	HashUUID       string
	VlessUUID      string
	TrojanPassword string
	SSPassword     string
}

type InboundUserData struct {
	Type string
	Tag  string

	Flow string

	CipherType CipherType
	IVCheck    bool
}

func BuildUserForInbound(inbound InboundUserData, user UserData) *protocol.User {
	const level uint32 = 0

	switch inbound.Type {
	case "vless":
		return BuildVlessUser(user.UserID, user.VlessUUID, inbound.Flow, level)
	case "trojan":
		return BuildTrojanUser(user.UserID, user.TrojanPassword, level)
	case "shadowsocks":
		return BuildShadowsocksUser(user.UserID, user.SSPassword, inbound.CipherType, inbound.IVCheck, level)
	default:
		return nil
	}
}

func ParseCipherType(s string) CipherType {
	switch s {
	case "aes-128-gcm", "AES_128_GCM":
		return CipherTypeAES128GCM
	case "aes-256-gcm", "AES_256_GCM":
		return CipherTypeAES256GCM
	case "chacha20-poly1305", "chacha20-ietf-poly1305", "CHACHA20_POLY1305":
		return CipherTypeCHACHA20POLY1305
	case "xchacha20-poly1305", "xchacha20-ietf-poly1305", "XCHACHA20_POLY1305":
		return CipherTypeXCHACHA20POLY1305
	case "none", "NONE":
		return CipherTypeNone
	default:
		return CipherTypeUnknown
	}
}
