package configs

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Logger struct {
	LogLevel string `yaml:"logLvl"`
	LogReq   bool   `yaml:"logReq"`
	LogResp  bool   `yaml:"logResp"`
}

type RemoteTrace struct {
	Enabled        bool    `yaml:"enabled"`
	TraceAgent     string  `yaml:"traceAgent"`
	TraceCollector string  `yaml:"traceCollector"`
	Ratio          float64 `yaml:"ratio"`
}

type RemoteProfiler struct {
	Enabled     bool   `yaml:"enabled"`
	ProfilerURL string `yaml:"profilerURL"`
}

type BaseConfig struct {
	Name           string          `yaml:"name"`
	Env            string          `yaml:"env"`
	Debug          bool            `yaml:"debug"`
	Logger         *Logger         `yaml:"logger"`
	Port           string          `yaml:"port"`
	RemoteTrace    *RemoteTrace    `yaml:"remoteTrace"`
	StatsEnabled   bool            `yaml:"statsEnabled"`
	RemoteProfiler *RemoteProfiler `yaml:"remoteProfiler"`
}

type Config struct {
	BaseConfig `yaml:",inline"`
	AppVersion string `yaml:"appVersion"`
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	rawContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(rawContent, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
