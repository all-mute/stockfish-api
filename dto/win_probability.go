package dto

type WinProbabilityRequest struct {
	Table string `json:"table"`
	Level uint16 `json:"level"`
}

type WinProbabilityResponse struct {
	Probability string `json:"probability"`
}

