package models

type GenerateSubtitleRequest struct {
	Model    string `json:"model"`
	Language string `json:"language"`
}

type GenerateSummaryRequest struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
}
