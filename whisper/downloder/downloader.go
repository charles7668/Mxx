package downloder

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	srcUrl   = "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/"
	modelExt = ".bin"
	bufSize  = 1024 * 64
)

var (
	DownloadingErr  = errors.New("model is already downloading")
	downloadingTask = make(map[string]context.CancelFunc)
	lock            sync.Mutex
)

// getTaskQueryKey generates a unique key for the download task based on the URL and output file path(absolute path)
func getTaskQueryKey(downloadUrl, outputFilePath string) string {
	return downloadUrl + ";" + outputFilePath
}

// Download downloads the model from the given URL and saves it to the specified path
func Download(ctx context.Context, modelName, dir string, progress func(float32, error)) error {
	if filepath.Ext(modelName) != modelExt {
		modelName += modelExt
	}
	downloadUrl := urlForModel(modelName)
	outputFile, _ := filepath.Abs(filepath.Join(dir, modelName))
	taskQueryKey := getTaskQueryKey(downloadUrl, outputFile)
	lock.Lock()
	if _, ok := downloadingTask[taskQueryKey]; ok {
		lock.Unlock()
		return DownloadingErr
	}

	// create context for download
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	downloadingTask[taskQueryKey] = cancelFunc
	lock.Unlock()
	go func() {
		client := &http.Client{}
		req, err := http.NewRequest("GET", downloadUrl, nil)
		if err != nil {
			progress(0, err)
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			progress(0, err)
			return
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		if resp.StatusCode != http.StatusOK {
			progress(0, errors.New("failed to download model: "+resp.Status))
			return
		}

		// If output file exists and is the same size as the model, skip
		if info, err := os.Stat(outputFile); err == nil && info.Size() == resp.ContentLength {
			progress(100, nil)
			return
		}

		// Create file
		w, err := os.Create(outputFile)
		if err != nil {
			progress(0, err)
			return
		}
		defer func(w *os.File) {
			_ = w.Close()
		}(w)

		// Progressively download the model
		data := make([]byte, bufSize)
		count, pct := int64(0), float32(0)
		ticker := time.NewTicker(5 * time.Second)
	ReceiveLoop:
		for {
			select {
			case <-cancelCtx.Done():
				progress(pct, cancelCtx.Err())
				break ReceiveLoop
			case <-ticker.C:
				pct = calculateProgressPercent(count, resp.ContentLength)
				progress(pct, nil)
			default:
				// Read body
				n, err := resp.Body.Read(data)
				if err != nil {
					pct = calculateProgressPercent(count, resp.ContentLength)
					progress(pct, err)
					break ReceiveLoop
				} else if m, err := w.Write(data[:n]); err != nil {
					progress(pct, err)
					break ReceiveLoop
				} else {
					count += int64(m)
					if count >= resp.ContentLength {
						progress(100, nil)
						break ReceiveLoop
					}
				}
			}
		}
		lock.Lock()
		delete(downloadingTask, taskQueryKey)
		lock.Unlock()
	}()

	return nil
}

// CancelDownload cancels the download of the model
func CancelDownload(modelName, dir string) {
	if filepath.Ext(modelName) != modelExt {
		modelName += modelExt
	}
	downloadUrl := urlForModel(modelName)
	outputFile, _ := filepath.Abs(filepath.Join(dir, modelName))
	if cancelFunc, ok := downloadingTask[getTaskQueryKey(downloadUrl, outputFile)]; ok {
		cancelFunc()
	}
}

// calculateProgressPercent calculates the progress percentage
func calculateProgressPercent(count, total int64) float32 {
	pct := float32(count) * float32(100.0) / float32(total)
	return pct
}

// urlForModel returns the URL for the given model
func urlForModel(model string) string {
	if !strings.HasPrefix(model, "ggml-") {
		model = "ggml-" + model
	}
	downloadUrl, _ := url.Parse(srcUrl)
	downloadUrl.Path = fmt.Sprintf("%s/%s", strings.TrimSuffix(downloadUrl.Path, "/"), model)
	return downloadUrl.String()
}
