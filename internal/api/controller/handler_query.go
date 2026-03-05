package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
