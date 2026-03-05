package controller

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
