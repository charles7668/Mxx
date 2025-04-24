package configs

import "os"

type ApiConfig struct {
	ModelStorePath string
	TempStorePath  string
	MediaStorePath string
}

var apiConfig *ApiConfig

// GetApiConfig get singleton instance of apiConfig
func GetApiConfig() *ApiConfig {
	if apiConfig == nil {
		apiConfig = &ApiConfig{
			ModelStorePath: "./data/models",
			TempStorePath:  "./data/temp",
			MediaStorePath: "./data/media",
		}
		err := os.RemoveAll(apiConfig.TempStorePath)
		if err != nil {
			panic(err)
		}
		err = os.RemoveAll(apiConfig.MediaStorePath)
		if err != nil {
			panic(err)
		}
		err = os.MkdirAll(apiConfig.TempStorePath, os.ModePerm)
		if err != nil {
			panic(err)
		}
		err = os.MkdirAll(apiConfig.ModelStorePath, os.ModePerm)
		if err != nil {
			panic(err)
		}
		err = os.MkdirAll(apiConfig.MediaStorePath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	return apiConfig
}
