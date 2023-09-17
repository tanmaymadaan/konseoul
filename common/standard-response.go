package common

type StandardResponse struct {
	StatusCode int                    `json:"statusCode"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"data"`
}
