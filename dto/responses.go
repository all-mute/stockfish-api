package dto

type MoveResponse struct {
	FenTable string `json:"fenTable"`
	Move     string `json:"move"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
