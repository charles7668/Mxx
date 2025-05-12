package models

type GenerateSubtitleRequest struct {
	Model    string `json:"model"`
	Language string `json:"language"`
}
