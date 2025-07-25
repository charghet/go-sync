package config

import (
	"os"
	"path/filepath"

	"github.com/charghet/go-sync/pkg/logger"
	"gopkg.in/yaml.v3"
)

type Config struct {
	User  UserConfig   `yaml:"user"`
	Repos []RepoConfig `yaml:"repos"`
}

type UserConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type RepoConfig struct {
	Path     string `yaml:"path"` // 本地路径
	Url      string `yaml:"url"`
	Branch   string `yaml:"branch"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Email    string `yaml:"email"`
	Debounce int    `yaml:"debounce"` // 防抖时间 秒
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
			Path:   repo.Path,
			Url:    repo.Url,
			Branch: repo.Branch,
		}
	}
	return res
}
