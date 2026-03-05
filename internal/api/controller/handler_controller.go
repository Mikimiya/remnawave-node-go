package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xtls/xray-core/features/inbound"

	"github.com/hteppl/remnawave-node-go/internal/logger"
	"github.com/hteppl/remnawave-node-go/internal/xray"
)

var (
	logFailedToParseAddUserRequest             = "Failed to parse add-user request"
	logFailedToParseAddUsersRequest            = "Failed to parse add-users request"
	logFailedToParseRemoveUserRequest          = "Failed to parse remove-user request"
	logFailedToParseRemoveUsersRequest         = "Failed to parse remove-users request"
	logFailedToParseGetInboundUsersRequest     = "Failed to parse get-inbound-users request"
	logFailedToParseGetInboundUsersCountReq    = "Failed to parse get-inbound-users-count request"
	logFailedToParseDropUsersConnectionsReq    = "Failed to parse drop-users-connections request"
	logFailedToParseDropIPsRequest             = "Failed to parse drop-ips request"
	logFailedToGetUserManager                  = "Failed to get user manager"
	logFailedToBuildUser                       = "Failed to build user - unsupported type"
	logFailedToAddUserToInbound                = "Failed to add user to inbound"
	logFailedToAddUserToInboundDuringBulk      = "Failed to add user to inbound during bulk add"
	logErrorRemovingUserFromInbounds           = "Error removing user from all inbounds (may not exist)"
	logErrorRemovingUserFromInboundsDuringBulk = "Error removing user from inbounds during bulk add"
	logErrorRemovingUserDuringBulkRemove       = "Error removing user from all inbounds during bulk remove"
	logUserAddedSuccessfully                   = "User added successfully"
	logBulkUsersAddedSuccessfully              = "Bulk users added successfully"
	logUserRemovedSuccessfully                 = "User removed successfully"
	logBulkUsersRemovedSuccessfully            = "Bulk users removed successfully"
)

type AddUserInboundData struct {
	Tag        string `json:"tag" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Type       string `json:"type" binding:"required"`
	UUID       string `json:"uuid,omitempty"`
	Flow       string `json:"flow,omitempty"`
	Password   string `json:"password,omitempty"`
	CipherType string `json:"cipherType,omitempty"`
	IVCheck    bool   `json:"ivCheck,omitempty"`
}

type AddUserHashData struct {
	VlessUUID     string `json:"vlessUuid,omitempty"`
	PrevVlessUUID string `json:"prevVlessUuid,omitempty"`
}

type AddUserRequest struct {
	Data     []AddUserInboundData `json:"data" binding:"required,dive"`
	HashData AddUserHashData      `json:"hashData"`
}

type AddUserResponseData struct {
	Success bool    `json:"success"`
	Error   *string `json:"error"`
}

type BulkUserData struct {
	UserID         string `json:"userId" binding:"required"`
	HashUUID       string `json:"hashUuid,omitempty"`
	VlessUUID      string `json:"vlessUuid,omitempty"`
	TrojanPassword string `json:"trojanPassword,omitempty"`
	SSPassword     string `json:"ssPassword,omitempty"`
}

type BulkInboundData struct {
	Tag        string `json:"tag" binding:"required"`
	Type       string `json:"type" binding:"required"`
	Flow       string `json:"flow,omitempty"`
	CipherType string `json:"cipherType,omitempty"`
	IVCheck    bool   `json:"ivCheck,omitempty"`
}

type BulkUserEntry struct {
	UserData    BulkUserData      `json:"userData" binding:"required"`
	InboundData []BulkInboundData `json:"inboundData" binding:"required,dive"`
}

type AddUsersRequest struct {
	AffectedInboundTags []string        `json:"affectedInboundTags"`
	Users               []BulkUserEntry `json:"users" binding:"required,dive"`
}

type RemoveUserHashData struct {
	VlessUUID string `json:"vlessUuid,omitempty"`
}

type RemoveUserRequest struct {
	Username string             `json:"username" binding:"required"`
	HashData RemoveUserHashData `json:"hashData"`
}

type BulkRemoveUserEntry struct {
	UserID   string `json:"userId" binding:"required"`
	HashUUID string `json:"hashUuid,omitempty"`
}

type RemoveUsersRequest struct {
	Users []BulkRemoveUserEntry `json:"users" binding:"required,dive"`
}

type GetInboundUsersRequest struct {
	Tag string `json:"tag" binding:"required"`
}

type InboundUser struct {
	Username string `json:"username"`
	Level    uint32 `json:"level"`
	Protocol string `json:"protocol"`
}

type GetInboundUsersResponseData struct {
	Users []InboundUser `json:"users"`
}

type GetInboundUsersCountRequest struct {
	Tag string `json:"tag" binding:"required"`
}

type GetInboundUsersCountResponseData struct {
	Count int64 `json:"count"`
}

type DropUsersConnectionsRequest struct {
	UserIDs []string `json:"userIds" binding:"required,min=1"`
}

type DropIPsRequest struct {
	IPs []string `json:"ips" binding:"required,min=1"`
}

type GenericResponseData struct {
	Success bool `json:"success"`
}

type HandlerController struct {
	core          *xray.Core
	configManager *xray.ConfigManager
	logger        *logger.Logger
}

func NewHandlerController(core *xray.Core, configManager *xray.ConfigManager, log *logger.Logger) *HandlerController {
	return &HandlerController{
		core:          core,
		configManager: configManager,
		logger:        log,
	}
}

func (c *HandlerController) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/add-user", c.handleAddUser)
	group.POST("/add-users", c.handleAddUsers)
	group.POST("/remove-user", c.handleRemoveUser)
	group.POST("/remove-users", c.handleRemoveUsers)
	group.POST("/get-inbound-users", c.handleGetInboundUsers)
	group.POST("/get-inbound-users-count", c.handleGetInboundUsersCount)
	group.POST("/drop-users-connections", c.handleDropUsersConnections)
	group.POST("/drop-ips", c.handleDropIPs)
}

func (c *HandlerController) getUserManager() (*xray.UserManager, error) {
	instance := c.core.Instance()
	if instance == nil {
		return nil, errors.New("xray core not running")
	}

	ibmFeature := instance.GetFeature(inbound.ManagerType())
	if ibmFeature == nil {
		return nil, errors.New("inbound manager not available")
	}

	ibm, ok := ibmFeature.(inbound.Manager)
	if !ok {
		return nil, errors.New("failed to cast to inbound manager")
	}

	return xray.NewUserManager(ibm, c.logger), nil
}

func (c *HandlerController) handleAddUser(ctx *gin.Context) {
	var req AddUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WithError(err).Error(logFailedToParseAddUserRequest)
		errMsg := err.Error()
		ctx.JSON(http.StatusOK, wrapResponse(AddUserResponseData{
			Success: false,
			Error:   &errMsg,
		}))
		return
	}

	if len(req.Data) == 0 {
		errMsg := "no inbound data provided"
		ctx.JSON(http.StatusOK, wrapResponse(AddUserResponseData{
			Success: false,
			Error:   &errMsg,
		}))
		return
	}

	userManager, err := c.getUserManager()
	if err != nil {
		c.logger.WithError(err).Error(logFailedToGetUserManager)
		errMsg := err.Error()
		ctx.JSON(http.StatusOK, wrapResponse(AddUserResponseData{
			Success: false,
			Error:   &errMsg,
		}))
		return
	}

	username := req.Data[0].Username
	bgCtx := context.Background()

	allTags := c.configManager.GetXtlsConfigInbounds()
	if err := userManager.RemoveUserFromAllInbounds(bgCtx, allTags, username); err != nil {
		c.logger.WithError(err).WithField("username", username).
			Warn(logErrorRemovingUserFromInbounds)
	}

	hashToRemove := req.HashData.PrevVlessUUID
	if hashToRemove == "" {
		hashToRemove = req.HashData.VlessUUID
	}
	if hashToRemove != "" {
		for _, tag := range allTags {
			c.configManager.RemoveUserFromInbound(tag, hashToRemove)
		}
	}

	for _, inboundData := range req.Data {
		userData := xray.UserData{
			UserID:    inboundData.Username,
			VlessUUID: inboundData.UUID,
		}

		if inboundData.Type == "trojan" {
			userData.TrojanPassword = inboundData.Password
		} else if inboundData.Type == "shadowsocks" {
			userData.SSPassword = inboundData.Password
		}

		inbound := xray.InboundUserData{
			Type:       inboundData.Type,
			Tag:        inboundData.Tag,
			Flow:       inboundData.Flow,
			CipherType: xray.ParseCipherType(inboundData.CipherType),
			IVCheck:    inboundData.IVCheck,
		}

		user := xray.BuildUserForInbound(inbound, userData)
		if user == nil {
			c.logger.WithField("type", inboundData.Type).
				WithField("tag", inboundData.Tag).
				Error(logFailedToBuildUser)
			continue
		}

		if err := userManager.AddUser(bgCtx, inboundData.Tag, user); err != nil {
			c.logger.WithError(err).
				WithField("tag", inboundData.Tag).
				WithField("username", inboundData.Username).
				Error(logFailedToAddUserToInbound)
			errMsg := err.Error()
			ctx.JSON(http.StatusOK, wrapResponse(AddUserResponseData{
				Success: false,
				Error:   &errMsg,
			}))
			return
		}
	}

	if req.HashData.VlessUUID != "" {
		for _, inboundData := range req.Data {
			c.configManager.AddUserToInbound(inboundData.Tag, req.HashData.VlessUUID)
		}
	}

	c.logger.WithField("username", username).
		WithField("inbounds", len(req.Data)).
		Info(logUserAddedSuccessfully)

	ctx.JSON(http.StatusOK, wrapResponse(AddUserResponseData{
		Success: true,
		Error:   nil,
	}))
}

func (c *HandlerController) handleAddUsers(ctx *gin.Context) {
	var req AddUsersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WithError(err).Error(logFailedToParseAddUsersRequest)
		errMsg := err.Error()
		ctx.JSON(http.StatusOK, wrapResponse(AddUserResponseData{
			Success: false,
			Error:   &errMsg,
		}))
		return
	}

	if len(req.Users) == 0 {
		ctx.JSON(http.StatusOK, wrapResponse(AddUserResponseData{
			Success: true,
			Error:   nil,
		}))
		return
	}

	userManager, err := c.getUserManager()
	if err != nil {
		c.logger.WithError(err).Error(logFailedToGetUserManager)
		errMsg := err.Error()
		ctx.JSON(http.StatusOK, wrapResponse(AddUserResponseData{
			Success: false,
			Error:   &errMsg,
		}))
		return
	}

	bgCtx := context.Background()

	allTags := req.AffectedInboundTags
	if len(allTags) == 0 {
		allTags = c.configManager.GetXtlsConfigInbounds()
	}

	for _, userEntry := range req.Users {
		username := userEntry.UserData.UserID
		hashUUID := userEntry.UserData.HashUUID

		if err := userManager.RemoveUserFromAllInbounds(bgCtx, allTags, username); err != nil {
			c.logger.WithError(err).WithField("username", username).
				Warn(logErrorRemovingUserFromInboundsDuringBulk)
		}

		if hashUUID != "" {
			for _, tag := range allTags {
				c.configManager.RemoveUserFromInbound(tag, hashUUID)
			}
		}

		for _, inboundData := range userEntry.InboundData {
			userData := xray.UserData{
				UserID:         username,
				HashUUID:       userEntry.UserData.HashUUID,
				VlessUUID:      userEntry.UserData.VlessUUID,
				TrojanPassword: userEntry.UserData.TrojanPassword,
				SSPassword:     userEntry.UserData.SSPassword,
			}

			inbound := xray.InboundUserData{
				Type:       inboundData.Type,
				Tag:        inboundData.Tag,
				Flow:       inboundData.Flow,
				CipherType: xray.ParseCipherType(inboundData.CipherType),
				IVCheck:    inboundData.IVCheck,
			}

			user := xray.BuildUserForInbound(inbound, userData)
			if user == nil {
				c.logger.WithField("type", inboundData.Type).
					WithField("tag", inboundData.Tag).
					Error(logFailedToBuildUser)
				continue
			}

			if err := userManager.AddUser(bgCtx, inboundData.Tag, user); err != nil {
				c.logger.WithError(err).
					WithField("tag", inboundData.Tag).
					WithField("username", username).
					Error(logFailedToAddUserToInboundDuringBulk)
				errMsg := err.Error()
				ctx.JSON(http.StatusOK, wrapResponse(AddUserResponseData{
					Success: false,
					Error:   &errMsg,
				}))
				return
			}

			if userEntry.UserData.HashUUID != "" {
				c.configManager.AddUserToInbound(inboundData.Tag, userEntry.UserData.HashUUID)
			}
		}
	}

	c.logger.WithField("count", len(req.Users)).Info(logBulkUsersAddedSuccessfully)

	ctx.JSON(http.StatusOK, wrapResponse(AddUserResponseData{
		Success: true,
		Error:   nil,
	}))
}

func (c *HandlerController) handleRemoveUser(ctx *gin.Context) {
	var req RemoveUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WithError(err).Error(logFailedToParseRemoveUserRequest)
		errMsg := err.Error()
		ctx.JSON(http.StatusOK, wrapResponse(AddUserResponseData{
			Success: false,
			Error:   &errMsg,
		}))
		return
	}

	userManager, err := c.getUserManager()
	if err != nil {
		c.logger.WithError(err).Error(logFailedToGetUserManager)
		errMsg := err.Error()
		ctx.JSON(http.StatusOK, wrapResponse(AddUserResponseData{
			Success: false,
			Error:   &errMsg,
		}))
		return
	}

	bgCtx := context.Background()

	allTags := c.configManager.GetXtlsConfigInbounds()
	if err := userManager.RemoveUserFromAllInbounds(bgCtx, allTags, req.Username); err != nil {
		c.logger.WithError(err).WithField("username", req.Username).
			Warn(logErrorRemovingUserFromInbounds)
	}

	if req.HashData.VlessUUID != "" {
		for _, tag := range allTags {
			c.configManager.RemoveUserFromInbound(tag, req.HashData.VlessUUID)
		}
	}

	c.logger.WithField("username", req.Username).Info(logUserRemovedSuccessfully)

	ctx.JSON(http.StatusOK, wrapResponse(AddUserResponseData{
		Success: true,
		Error:   nil,
	}))
}

func (c *HandlerController) handleRemoveUsers(ctx *gin.Context) {
	var req RemoveUsersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WithError(err).Error(logFailedToParseRemoveUsersRequest)
		errMsg := err.Error()
		ctx.JSON(http.StatusOK, wrapResponse(AddUserResponseData{
			Success: false,
			Error:   &errMsg,
		}))
		return
	}

	if len(req.Users) == 0 {
		ctx.JSON(http.StatusOK, wrapResponse(AddUserResponseData{
			Success: true,
			Error:   nil,
		}))
		return
	}

	userManager, err := c.getUserManager()
	if err != nil {
		c.logger.WithError(err).Error(logFailedToGetUserManager)
		errMsg := err.Error()
		ctx.JSON(http.StatusOK, wrapResponse(AddUserResponseData{
			Success: false,
			Error:   &errMsg,
		}))
		return
	}

	bgCtx := context.Background()
	allTags := c.configManager.GetXtlsConfigInbounds()

	for _, userEntry := range req.Users {
		if err := userManager.RemoveUserFromAllInbounds(bgCtx, allTags, userEntry.UserID); err != nil {
			c.logger.WithError(err).WithField("username", userEntry.UserID).
				Warn(logErrorRemovingUserDuringBulkRemove)
		}

		if userEntry.HashUUID != "" {
			for _, tag := range allTags {
				c.configManager.RemoveUserFromInbound(tag, userEntry.HashUUID)
			}
		}
	}

	c.logger.WithField("count", len(req.Users)).Info(logBulkUsersRemovedSuccessfully)

	ctx.JSON(http.StatusOK, wrapResponse(AddUserResponseData{
		Success: true,
		Error:   nil,
	}))
}

func (c *HandlerController) handleGetInboundUsers(ctx *gin.Context) {
	var req GetInboundUsersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WithError(err).Error(logFailedToParseGetInboundUsersRequest)
		ctx.JSON(http.StatusOK, wrapResponse(GetInboundUsersResponseData{
			Users: []InboundUser{},
		}))
		return
	}

	userManager, err := c.getUserManager()
	if err != nil {
		ctx.JSON(http.StatusOK, wrapResponse(GetInboundUsersResponseData{
			Users: []InboundUser{},
		}))
		return
	}

	bgCtx := context.Background()
	users, err := userManager.GetInboundUsers(bgCtx, req.Tag)
	if err != nil {
		ctx.JSON(http.StatusOK, wrapResponse(GetInboundUsersResponseData{
			Users: []InboundUser{},
		}))
		return
	}

	inboundUsers := make([]InboundUser, 0, len(users))
	for _, u := range users {
		inboundUsers = append(inboundUsers, InboundUser{
			Username: u.Email,
			Level:    u.Level,
		})
	}

	ctx.JSON(http.StatusOK, wrapResponse(GetInboundUsersResponseData{
		Users: inboundUsers,
	}))
}

func (c *HandlerController) handleGetInboundUsersCount(ctx *gin.Context) {
	var req GetInboundUsersCountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WithError(err).Error(logFailedToParseGetInboundUsersCountReq)
		ctx.JSON(http.StatusOK, wrapResponse(GetInboundUsersCountResponseData{
			Count: 0,
		}))
		return
	}

	userManager, err := c.getUserManager()
	if err != nil {
		ctx.JSON(http.StatusOK, wrapResponse(GetInboundUsersCountResponseData{
			Count: 0,
		}))
		return
	}

	bgCtx := context.Background()
	count, err := userManager.GetInboundUsersCount(bgCtx, req.Tag)
	if err != nil {
		ctx.JSON(http.StatusOK, wrapResponse(GetInboundUsersCountResponseData{
			Count: 0,
		}))
		return
	}

	ctx.JSON(http.StatusOK, wrapResponse(GetInboundUsersCountResponseData{
		Count: count,
	}))
}

func (c *HandlerController) handleDropUsersConnections(ctx *gin.Context) {
	var req DropUsersConnectionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WithError(err).Error(logFailedToParseDropUsersConnectionsReq)
		ctx.JSON(http.StatusOK, wrapResponse(GenericResponseData{
			Success: false,
		}))
		return
	}

	ctx.JSON(http.StatusOK, wrapResponse(GenericResponseData{
		Success: true,
	}))
}

func (c *HandlerController) handleDropIPs(ctx *gin.Context) {
	var req DropIPsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.WithError(err).Error(logFailedToParseDropIPsRequest)
		ctx.JSON(http.StatusOK, wrapResponse(GenericResponseData{
			Success: false,
		}))
		return
	}

	ctx.JSON(http.StatusOK, wrapResponse(GenericResponseData{
		Success: true,
	}))
}
