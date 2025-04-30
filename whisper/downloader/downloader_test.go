package downloader

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDownload(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()
	modelName := "tiny.en"
	filePath := filepath.Join(tempDir, modelName+".bin")
	removeFileAfter := func() {
		_ = os.Remove(filePath)
	}
	defer removeFileAfter()

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	err := Download(ctx, modelName, tempDir, func(progress float32, err error) {
		if err != nil {
			t.Errorf("error while download : %v", err)
			cancel()
			return
		}

		if progress >= 100 {
			t.Logf("Download Progress: %.2f%%", progress)
			// check if the file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Fatalf("download file not exist: %v", err)
			}
			cancel()
		} else {
			t.Logf("Download Progress: %.2f%%", progress)
		}
	})

	// test model already exists
	if err != nil {
		t.Fatalf("download failed: %v", err)
	}
	<-ctx.Done()
	ctx, cancel = context.WithCancel(context.Background())
	var downloadErr error
	err = Download(ctx, modelName, tempDir, func(progress float32, err error) {
		defer cancel()
		if progress < 100 {
			downloadErr = errors.New("when file already exists, progress should be 100")
		}
	})
	if err != nil {
		t.Fatalf("download failed: %v", err)
	}
	<-ctx.Done()
	if downloadErr != nil {
		t.Fatalf("download failed: %v", downloadErr)
	}
}

func TestDownloadingCheck(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()
	modelName := "tiny.en"
	filePath := filepath.Join(tempDir, modelName+".bin")
	removeFileAfter := func() {
		_ = os.Remove(filePath)
	}
	defer removeFileAfter()

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// start multiple download tasks , only one should be allowed to start
	for i := 0; i < 100; i++ {
		go func() {
			err := Download(ctx, modelName, tempDir, func(progress float32, err error) {
				if err != nil {
					t.Errorf("%d error while download : %v", i, err)
					cancel()
					return
				}

				if progress >= 100 {
					t.Logf("%d Download Progress: %.2f%%", i, progress)
					// check if the file exists
					if _, err := os.Stat(filePath); os.IsNotExist(err) {
						t.Errorf("download file not exist: %v", err)
						cancel()
						return
					}
					cancel()
				} else {
					t.Logf("%d Download Progress: %.2f%%", i, progress)
				}
			})
			if err != nil && !errors.Is(err, AlreadyDownloadedErr) {
				t.Errorf("%d download failed: %v", i, err)
				cancel()
				return
			}
		}()
	}

	<-ctx.Done()
}

func TestCancelDownload(t *testing.T) {
	tempDir := t.TempDir()
	modelName := "tiny.en"
	filePath := filepath.Join(tempDir, modelName+".bin")
	removeFileAfter := func() {
		_ = os.Remove(filePath)
	}
	defer removeFileAfter()

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	err := Download(ctx, modelName, tempDir, func(progress float32, err error) {
		if err != nil {
			if err.Error() == "context canceled" {
				t.Logf("Download canceled: %v", err)
				cancel()
				return
			}
			t.Errorf("error while download : %v", err)
			cancel()
			return
		}

		if progress >= 100 {
			t.Logf("Download Progress: %.2f%%", progress)
			// check if the file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Fatalf("download file not exist: %v", err)
			}
			cancel()
		} else {
			t.Logf("Download Progress: %.2f%%", progress)
		}
	})

	if err != nil {
		t.Fatalf("download failed: %v", err)
	}

	// cancel after 5 seconds
	time.Sleep(1 * time.Second)
	CancelDownload(modelName, filePath)

	<-ctx.Done()
}
