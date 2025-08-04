package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/charghet/go-sync/pkg/logger"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig `yaml:"server"`
	User     UserConfig   `yaml:"user"`
	Repos    []RepoConfig `yaml:"repos"`
	Ignore   *int         `yaml:"ignore"`
	Pull     *bool        `yaml:"pull"`
	Debounce *int         `yaml:"debounce" json:"debounce"`
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
	Ignore   *int   `yaml:"ignore"`
	Pull     *bool  `yaml:"pull"`
	Debounce *int   `yaml:"debounce" json:"debounce"` // 防抖时间 秒
}

var path = "config.yaml"
var config *Config

func SetDefaultConfig(con *Config) {
	if con.Ignore == nil {
		i := 3
		con.Ignore = &i
	}

	if con.Server.Host == "" {
		con.Server.Host = "127.0.0.1"
	}
	if con.Server.Port == 0 {
		con.Server.Port = 2222
	}

	if con.User.Username == "" {
		con.User.Username = "admin"
	}
	if con.User.Password == "" {
		con.User.Password = "admin123"
	}

	for i := range con.Repos {
		r := &con.Repos[i]
		if r.Name == "" {
			r.Name = r.Path
		}
		if r.Branch == "" {
			r.Branch = "master"
		}
		if r.Email == "" {
			r.Email = "go-sync@example.com"
		}

		if r.Ignore == nil {
			if con.Ignore == nil {
				i := 3
				r.Ignore = &i
			} else {
				r.Ignore = con.Ignore
			}
		}

		if r.Debounce == nil {
			if con.Debounce == nil {
				i := 3
				r.Debounce = &i
			} else {
				r.Debounce = con.Debounce
			}
		}

		if r.Pull == nil {
			if con.Pull == nil {
				p := true
				r.Pull = &p
			} else {
				r.Pull = con.Pull
			}
		}
	}
}

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

	SetDefaultConfig(config)
	b, _ := json.Marshal(config)
	logger.Debug("config:", string(b))
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
