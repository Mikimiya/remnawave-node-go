package controller

type successResponse struct {
	Response interface{} `json:"response"`
}

func wrapResponse(data interface{}) successResponse {
	return successResponse{Response: data}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
