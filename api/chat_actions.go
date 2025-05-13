package api

import (
	"Mxx/api/constant"
	"Mxx/api/graceful"
	"Mxx/api/models"
	"Mxx/api/subtitle"
	"Mxx/api/task"
	"Mxx/llm"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSummary(c *gin.Context) {
	var requestBody models.GenerateSummaryRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request body"})
		return
	}
	logger := getLoggerFromContext(c)
	sessionId := c.GetString(constant.SessionIdCtxKey)
	subtitleManager := subtitle.GetManager()
	if subtitleManager == nil {
		logger.Error("subtitle manager is nil")
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "subtitle not found"})
		return
	}
	if !subtitleManager.Exist(sessionId) {
		logger.Error("subtitle not exist")
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "subtitle not exist"})
		return
	}
	if state, found := task.GetTaskState(sessionId); found && state.Status == task.Running {
		logger.Sugar().Errorf("another task is running task : %s", state.Task)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "another task is running"})
		return
	}
	task.StartTask(sessionId, task.State{
		Task: "generating summary",
	})
	defer task.CompleteTask(sessionId)
	subtitleString := subtitleManager.ToPlainText(sessionId)

	provider, err := llm.PrepareRunner(toProviderEnum(requestBody.Provider), map[string]string{
		"model": requestBody.Model,
	})
	if err != nil {
		logger.Sugar().Errorf("prepare runner error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: fmt.Sprintf("prepare %s provider error", requestBody.Provider)})
		task.FailedTask(sessionId, err)
		return
	}
	ctx := graceful.BackgroundContext
	chatResult, err := provider.Chat(ctx, "You are an excellent analyst. Please help me generate a summary based on the following video subtitles.",
		subtitleString, "")
	if err != nil {
		logger.Sugar().Errorf("Chat() returned an error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Chat() returned an error"})
		task.FailedTask(sessionId, err)
		return
	}
	c.JSON(http.StatusOK, models.SummaryResponse{
		Status:  http.StatusOK,
		Summary: chatResult,
	})
}

func toProviderEnum(providerString string) llm.Provider {
	switch providerString {
	case "ollama":
		return llm.Ollama
	}
	return llm.Unknown
}
