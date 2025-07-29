package config

import (
	"os"
	"path/filepath"

	"github.com/charghet/go-sync/pkg/logger"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	User   UserConfig   `yaml:"user"`
	Repos  []RepoConfig `yaml:"repos"`
	Ignore int          `yaml:"ignore"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type UserConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type RepoConfig struct {
	Name     string `yaml:"name" json:"name"`
	Path     string `yaml:"path" json:"path"` // 本地路径
	Url      string `yaml:"url" json:"url"`
	Branch   string `yaml:"branch" json:"branch"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Email    string `yaml:"email" json:"email"`
	Debounce int    `yaml:"debounce" json:"debounce"` // 防抖时间 秒
}

var path = "config.yaml"
var config *Config

func toInit() {
	if config != nil {
		return
	}
	config = &Config{}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		logger.Fatal("Failed to read config file:", err)
	}
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		logger.Fatal("Failed to decode config:", err)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		logger.Warn("Failed to get absolute path of config file:", err)
	}
	logger.Info("loaded config:", abs)
}

func GetConfig() *Config {
	toInit()
	return config
}

func SetPath(p string) *Config {
	path = p
	toInit()
	return config
}

func RepoInfo() []RepoConfig {
	res := make([]RepoConfig, len(config.Repos))
	for i, repo := range config.Repos {
		res[i] = RepoConfig{
			Name:   repo.Name,
			Path:   repo.Path,
			Url:    repo.Url,
			Branch: repo.Branch,
		}
	}
	return res
}
