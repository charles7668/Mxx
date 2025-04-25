package api

import (
	"Mxx/api/configs"
	"Mxx/api/graceful"
	"Mxx/api/media"
	"Mxx/api/session"
	"Mxx/api/task"
	"Mxx/ffmpeg/converter"
	"Mxx/whisper/downloder"
	"Mxx/whisper/transcription"
	"context"
	"errors"
	"fmt"
	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func generateSessionId(c *gin.Context) {
	sessionId := session.GenerateSessionId()
	session.AddToManager(sessionId, time.Now())
	c.JSON(200, gin.H{"session_id": sessionId})
}

func mediaUpload(c *gin.Context) {
	sessionId := c.GetHeader("X-Session-Id")
	storeDir := filepath.Join(configs.GetApiConfig().MediaStorePath, sessionId)
	err := os.MkdirAll(storeDir, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}
	if stat, err := os.Stat(storeDir); err != nil || !stat.IsDir() {
		// If the sessionId is not a directory, return an error because a required file might have the same name as this ID.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is invalid"})
		return
	}

	if state, found := task.GetTaskState(sessionId); found && state.Status == task.Running {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Other task is running Task : " + state.Task})
		return
	}

	task.StartTask(sessionId, task.State{
		Task: "uploading files",
	})
	defer task.CompleteTask(sessionId)

	// get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	targetPath := filepath.Join(storeDir, file.Filename)
	if err := c.SaveUploadedFile(file, targetPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	mediaManager := media.GetMediaManager()
	mediaManager.SetMediaPath(sessionId, targetPath)
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file_path": targetPath})
}

func getUploadedMedia(c *gin.Context) {
	sessionId := c.GetHeader("X-Session-Id")
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
	sessionId := c.GetHeader("X-Session-Id")
	mediaManager := media.GetMediaManager()
	mediaPath := mediaManager.GetMediaPath(sessionId)
	if mediaPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No media file found"})
		return
	}
	if state, found := task.GetTaskState(sessionId); found && state.Status == task.Running {
		c.JSON(http.StatusBadRequest, gin.H{
			"is_running": true,
			"error":      "Other task is running Task : " + state.Task,
		})
		return
	}
	task.StartTask(sessionId, task.State{
		Task: "generating subtitles",
	})
	go func() {
		task.StartTask(sessionId, task.State{
			Task: "downloading model",
		})
		apiConfig := configs.GetApiConfig()
		modelPath := apiConfig.ModelStorePath
		downloadCtx, downloadCancel := context.WithCancel(context.Background())
		var downloadErr error = nil
		err := downloder.Download(downloadCtx, "tiny", modelPath, func(progress float32, err error) {
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
			task.FailedTask(sessionId, downloadErr)
			return
		}
		<-downloadCtx.Done()
		if downloadErr != nil {
			task.FailedTask(sessionId, downloadErr)
			return
		}
		task.StartTask(sessionId, task.State{
			Task: "start convert file to wav",
		})
		audioConverter := converter.CreateAudioConverter("ffmpeg")
		mediaManager = media.GetMediaManager()
		inputFilePath := mediaManager.GetMediaPath(sessionId)
		tempDir := apiConfig.TempStorePath
		tempUUID := filepath.Join(tempDir, sessionId)
		err = os.MkdirAll(tempUUID, os.ModePerm)
		if err != nil {
			task.FailedTask(sessionId, err)
			return
		}
		audioTarget := filepath.Join(tempUUID, "output.wav")
		err = audioConverter.Convert(inputFilePath, audioTarget)
		if err != nil {
			task.FailedTask(sessionId, err)
			return
		}
		task.StartTask(sessionId, task.State{
			Task: "start generate subtitles",
		})
		// write subtitles to file
		subtitleFile := filepath.Join(tempUUID, "output.txt")
		stream, err := os.Create(subtitleFile)
		if err != nil {
			task.FailedTask(sessionId, err)
			return
		}
		defer func(stream *os.File) {
			err := stream.Close()
			if err != nil {
				fmt.Println("failed to close file: ", err)
			}
		}(stream)

		var whisperErr error = nil
		whisperContext, whisperCancelFunc := context.WithCancel(graceful.BackgroundContext)
		whisperOptions := transcription.CreateOptions()
		whisperOptions.SegmentCallback = func(segment whisper.Segment) {
			writeString := fmt.Sprintf("[%6s -> %6s] %s", segment.Start.Truncate(time.Millisecond), segment.End.Truncate(time.Millisecond), segment.Text)
			_, writeErr := stream.WriteString(writeString + "\n")
			if writeErr != nil {
				whisperCancelFunc()
				whisperErr = writeErr
			}
		}
		err = transcription.Transcribe(whisperContext, audioTarget, filepath.Join(modelPath, "tiny"), whisperOptions)
		if err != nil {
			whisperCancelFunc()
			task.FailedTask(sessionId, err)
			return
		}
		if whisperErr != nil {
			whisperCancelFunc()
			task.FailedTask(sessionId, whisperErr)
			return
		}
		task.CompleteTask(sessionId)
	}()
	c.JSON(200, gin.H{"message": "Generate subtitles task started"})
}

func getMediaTaskState(c *gin.Context) {
	sessionId := c.GetHeader("X-Session-Id")
	state, found := task.GetTaskState(sessionId)
	if !found {
		state.Status = task.Completed
	}
	status := state.String()
	c.JSON(http.StatusOK, gin.H{
		"task":   state.Task,
		"status": status,
	})
}

func getSubtitle(c *gin.Context) {
	sessionId := c.GetHeader("X-Session-Id")
	tempDir := configs.GetApiConfig().TempStorePath
	subTitleFile := filepath.Join(tempDir, sessionId, "output.txt")
	if stat, err := os.Stat(subTitleFile); errors.Is(err, os.ErrNotExist) || stat.IsDir() {
		c.JSON(http.StatusNotFound, gin.H{"error": "No subtitle found"})
		return
	}
	file, err := os.Open(subTitleFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open subtitle file"})
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("failed to close file: ", err)
		}
	}(file)
	content, err := os.ReadFile(subTitleFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read subtitle file"})
		return
	}
	c.Header("Content-Type", "text/plain")
	c.JSON(http.StatusOK, gin.H{
		"result": string(content),
	})
}
