package types

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ResponseWithData struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Success   bool   `json:"success"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

type ErrorResponseWithData struct {
	Success   bool        `json:"success"`
	ErrorCode string      `json:"error_code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}
