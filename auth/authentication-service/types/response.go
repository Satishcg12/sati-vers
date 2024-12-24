package types

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ResponseWithDetails struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
