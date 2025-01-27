package main

type Config struct {
	DefaultPath     string
	MaxFileSize     int64
	ExcludePatterns []string
	IncludeHidden   bool
	OutputFormat    string
}

func NewDefaultConfig() *Config {
	return &Config{
		DefaultPath:     ".",
		MaxFileSize:     1024 * 1024 * 100, // 100MB
		ExcludePatterns: []string{".git", "node_modules"},
		IncludeHidden:   true,
		OutputFormat:    "xml",
	}
}
