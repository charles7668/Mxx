package api

import (
	"Mxx/api/configs"
	"Mxx/api/constant"
	"Mxx/api/graceful"
	"Mxx/api/log"
	"Mxx/api/media"
	"Mxx/api/models"
	"Mxx/api/session"
	"Mxx/api/subtitle"
	"Mxx/api/task"
	"Mxx/ffmpeg/converter"
	"Mxx/whisper/downloader"
	"Mxx/whisper/transcription"
	"context"
	"errors"
	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func generateSessionId(c *gin.Context) {
	sessionId := session.GenerateSessionId()
	generateTime := time.Now()
	log.GetLogger().Sugar().Infof("generate session id %s at %s", sessionId, generateTime.Format(time.RFC3339))
	session.Update(sessionId, generateTime)
	c.JSON(http.StatusOK, &models.SessionResponse{
		Status:    http.StatusOK,
		SessionId: sessionId,
	})
}

func mediaUpload(c *gin.Context) {
	sessionId := c.GetString(constant.SessionIdCtxKey)
	storeDir := filepath.Join(configs.GetApiConfig().MediaStorePath, sessionId)
	err := os.MkdirAll(storeDir, os.ModePerm)
	logger := log.GetLogger()
	if err != nil {
		logger.Sugar().Errorf("failed to create directory for session : %s , err : %s", sessionId, err.Error())
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "failed to create directory",
		})
		return
	}

	// If the sessionId is not a directory, return an error because a required file might have the same name as this ID.
	if stat, err := os.Stat(storeDir); err != nil || !stat.IsDir() {
		logger.Sugar().Errorf("sessionId is invalid : %s", sessionId)
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status: http.StatusUnauthorized,
			Error:  "Session ID is invalid",
		})
		return
	}

	if state, found := task.GetTaskState(sessionId); found && state.Status == task.Running {
		logger.Sugar().Infof("Another task is running : %s", state.Task)
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Other task is running Task : " + state.Task,
		})
		return
	}

	logger.Info("start upload file")
	defer logger.Info("finish upload file")
	task.StartTask(sessionId, task.State{
		Task: "uploading files",
	})
	defer task.CompleteTask(sessionId)

	// get file from form
	file, err := c.FormFile("file")
	if err != nil {
		logger.Error("failed to get file from form")
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "File is required",
		})
		return
	}

	targetPath := filepath.Join(storeDir, file.Filename)
	logger.Sugar().Infof("try save file to %s", targetPath)
	if err := c.SaveUploadedFile(file, targetPath); err != nil {
		logger.Sugar().Errorf("failed to save file : %s , err : %s", targetPath, err.Error())
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failed to save file",
		})
		return
	}
	mediaManager := media.GetMediaManager()
	mediaManager.SetMediaPath(sessionId, targetPath)

	m3u8Converter := converter.CreateM3U8Converter("")
	m3u8Target := filepath.Join(storeDir, "output.m3u8")
	err = m3u8Converter.Convert(targetPath, m3u8Target)
	if err != nil {
		logger.Sugar().Errorf("failed to save m3u8 : %s , err : %s", targetPath, err.Error())
		c.JSON(http.StatusInternalServerError, &models.ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failed to convert file to m3u8",
		})
		return
	}

	c.JSON(http.StatusOK, &models.FileUploadResponse{
		Status: http.StatusOK,
		File:   file.Filename,
	})
}

func getUploadedMedia(c *gin.Context) {
	obj, _ := c.Get(constant.SessionIdCtxKey)
	sessionId := obj.(string)
	mediaManager := media.GetMediaManager()
	mediaPath := mediaManager.GetMediaPath(sessionId)
	if mediaPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No media file found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"file_name": filepath.Base(mediaPath),
	})
}

func generateMediaSubtitles(c *gin.Context) {
	sessionId := c.GetString(constant.SessionIdCtxKey)
	var body models.GenerateSubtitleRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		log.GetLogger().Sugar().Errorf("failed to bind json : %s", err.Error())
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Invalid request",
		})
		return
	}
	mediaManager := media.GetMediaManager()
	mediaPath := mediaManager.GetMediaPath(sessionId)
	logger := log.GetLogger()
	if mediaPath == "" {
		logger.Sugar().Errorf("No media file found for session ID: %s", sessionId)
		c.JSON(http.StatusNotFound, &models.ErrorResponse{
			Status: http.StatusNotFound,
			Error:  "No media file found",
		})
		return
	}
	if state, found := task.GetTaskState(sessionId); found && state.Status == task.Running {
		logger.Sugar().Infof("Another task is running : %s", state.Task)
		c.JSON(http.StatusBadRequest, &models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Another task is running : " + state.Task,
		})
		return
	}
	logger.Info("start generate media subtitles")
	task.StartTask(sessionId, task.State{
		Task: "start generate media subtitles",
	})
	go func() {
		defer logger.Info("finish generate media subtitles")
		logger.Info("start download model")
		task.StartTask(sessionId, task.State{
			Task: "downloading model",
		})
		defer task.CompleteTask(sessionId)
		apiConfig := configs.GetApiConfig()
		modelPath := apiConfig.ModelStorePath
		downloadCtx, downloadCancel := context.WithCancel(context.Background())
		var downloadErr error = nil
		err := downloader.Download(downloadCtx, body.Model, modelPath, func(progress float32, err error) {
			if err != nil {
				downloadCancel()
				downloadErr = err
				return
			}
			if progress >= 100 {
				downloadCancel()
			}
		})
		if err != nil {
			if errors.Is(err, downloader.AlreadyDownloadedErr) {
				logger.Sugar().Infof("model already downloaded")
				downloadCancel()
			} else {
				logger.Sugar().Errorf("download error : %s", err.Error())
				task.FailedTask(sessionId, downloadErr)
				return
			}
		}
		<-downloadCtx.Done()
		if downloadErr != nil {
			logger.Sugar().Errorf("download error : %s", err.Error())
			task.FailedTask(sessionId, downloadErr)
			return
		}
		logger.Info("start convert file to wav")
		task.StartTask(sessionId, task.State{
			Task: "converting file to wav",
		})
		audioConverter := converter.CreateAudioConverter("")
		mediaManager = media.GetMediaManager()
		inputFilePath := mediaManager.GetMediaPath(sessionId)
		tempDir := apiConfig.TempStorePath
		tempUUID := filepath.Join(tempDir, sessionId)
		err = os.MkdirAll(tempUUID, os.ModePerm)
		if err != nil {
			logger.Sugar().Errorf("failed to create temp dir: %s, err: %s", tempUUID, err)
			task.FailedTask(sessionId, err)
			return
		}
		audioTarget := filepath.Join(tempUUID, "output.wav")
		err = audioConverter.Convert(inputFilePath, audioTarget)
		if err != nil {
			logger.Sugar().Errorf("failed to convert file: %s, err: %s", inputFilePath, err)
			task.FailedTask(sessionId, err)
			return
		}
		logger.Info("start generate subtitles")
		task.StartTask(sessionId, task.State{
			Task: "generating subtitles",
		})
		// write subtitles to file
		logger.Info("start generate subtitles by whisper")
		subtitleManager := subtitle.GetManager()
		subtitleManager.Clear(sessionId)
		var whisperErr error = nil
		whisperContext, whisperCancelFunc := context.WithCancel(graceful.BackgroundContext)
		whisperOptions := transcription.CreateOptions()
		whisperOptions.SegmentCallback = func(segment whisper.Segment) {
			store := subtitle.Segment{
				StartTime: segment.Start,
				EndTime:   segment.End,
				Text:      segment.Text,
			}
			subtitleManager.Add(sessionId, store)
		}
		err = transcription.Transcribe(whisperContext, audioTarget, filepath.Join(modelPath, body.Model), whisperOptions)
		if err != nil {
			whisperCancelFunc()
			logger.Sugar().Errorf("failed to transcribe file: %s, err: %s", audioTarget, err)
			task.FailedTask(sessionId, err)
			return
		}
		whisperCancelFunc()
		if whisperErr != nil {
			logger.Sugar().Errorf("failed to transcribe file: %s, err: %s", audioTarget, whisperErr)
			task.FailedTask(sessionId, whisperErr)
			return
		}
	}()
	c.JSON(200, struct{}{})
}

func getMediaTaskState(c *gin.Context) {
	sessionId := c.GetString(constant.SessionIdCtxKey)
	state, found := task.GetTaskState(sessionId)
	if !found {
		state.Status = task.Completed
	}
	status := state.String()
	c.JSON(http.StatusOK, &models.TaskStateResponse{
		Status:    http.StatusOK,
		Task:      state.Task,
		TaskState: status,
	})
}

func getSubtitle(c *gin.Context) {
	sessionId := c.GetString(constant.SessionIdCtxKey)
	subtitleManager := subtitle.GetManager()
	if !subtitleManager.Exist(sessionId) {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Status: http.StatusNotFound,
			Error:  "subtitle not found",
		})
		return
	}
	var builder strings.Builder

	for _, segment := range subtitleManager.GetSegments(sessionId) {
		builder.WriteString("[")
		builder.WriteString(segment.StartTime.String())
		builder.WriteString(" -> ")
		builder.WriteString(segment.EndTime.String())
		builder.WriteString("] ")
		builder.WriteString(segment.Text)
		builder.WriteString("\n")
	}
	content := []byte(builder.String())

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, models.ValueResponse{
		Status: http.StatusOK,
		Value:  string(content),
	})
}

func getASSFormatSubtitle(c *gin.Context) {
	sessionId := c.GetString(constant.SessionIdCtxKey)
	subtitleManager := subtitle.GetManager()
	if !subtitleManager.Exist(sessionId) {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Status: http.StatusNotFound,
			Error:  "subtitle not found",
		})
		return
	}
	assContent := subtitleManager.ToASS(sessionId)
	c.Header("Content-Disposition", "attachment; filename=subtitles.ass")
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, assContent)
}

func getPreviewMediaList(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Token is required",
		})
		return
	}
	mediaManager := media.GetMediaManager()
	mediaPath := mediaManager.GetMediaPath(token)
	mediaDir := path.Dir(mediaPath)
	m3u8Path := filepath.Join(mediaDir, "output.m3u8")
	if _, err := os.Stat(m3u8Path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Status: http.StatusNotFound,
			Error:  "No media file found",
		})
		return
	}
	c.Header("Content-Type", "application/vnd.apple.mpegurl")
	c.File(m3u8Path)
}

func getPreviewMediaFile(c *gin.Context) {
	token := c.Param("token")
	requestFile := c.Param("segment")
	if token == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Token is required",
		})
		return
	}
	mediaManager := media.GetMediaManager()
	mediaPath := mediaManager.GetMediaPath(token)
	mediaDir := path.Dir(mediaPath)
	tsPath := filepath.Join(mediaDir, requestFile)
	if _, err := os.Stat(tsPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Status: http.StatusNotFound,
			Error:  "No media file found",
		})
		return
	}
	c.Header("Content-Type", "application/vnd.apple.mpegurl")
	c.File(tsPath)
}
