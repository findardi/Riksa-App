package config

import "time"

type ViewerConfig struct {
	GotenbergURL      string
	ConvertTimeout    time.Duration
	DPI               int
	RenderTimeout     time.Duration
	RenderConcurrency int
}

func LoadViewerConfig() (ViewerConfig, error) {
	convertTimeout, err := GetEnvDuration("VIEWER_CONVERT_TIMEOUT", 2*time.Minute)
	if err != nil {
		return ViewerConfig{}, err
	}

	renderTimeout, err := GetEnvDuration("VIEWER_RENDER_TIMEOUT", 30*time.Second)
	if err != nil {
		return ViewerConfig{}, err
	}

	dpi, err := GetEnvInt("VIEWER_DPI", 150)
	if err != nil {
		return ViewerConfig{}, err
	}

	concurrency, err := GetEnvInt("VIEWER_RENDER_CONCURRENCY", 4)
	if err != nil {
		return ViewerConfig{}, err
	}

	if dpi <= 0 {
		dpi = 150
	}

	if concurrency <= 0 {
		concurrency = 1
	}

	return ViewerConfig{
		GotenbergURL:      GetEnv("GOTENBERG_URL", "http://localhost:3000"),
		ConvertTimeout:    convertTimeout,
		DPI:               dpi,
		RenderTimeout:     renderTimeout,
		RenderConcurrency: concurrency,
	}, nil
}
