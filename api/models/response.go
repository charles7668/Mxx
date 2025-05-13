package models

type ValueResponse struct {
	Status int `json:"status"`
	Value  any `json:"value"`
}

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

type SessionResponse struct {
	Status    int    `json:"status"`
	SessionId string `json:"sessionId"`
}

type FileUploadResponse struct {
	Status int    `json:"status"`
	File   string `json:"file"`
}

type TaskStateResponse struct {
	Status    int    `json:"status"`
	Task      string `json:"task"`
	TaskState string `json:"taskState"`
}

type SummaryResponse struct {
	Status  int    `json:"status"`
	Summary string `json:"summary"`
}
