package configs

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"os"
	"security-webhook/utils/log"
)

const DefaultConfigPath = "configs/default.yaml"

var GlobalConfig *Config

func init() {
	if err := loadConfigFromFile(DefaultConfigPath, &GlobalConfig); err != nil {
		log.Logger.Info("load default config failed", zap.Error(err))
		GlobalConfig.CheckItems = CheckItems{
			ForbiddenPrivilegedContainer: true,
		}
	}

	log.Logger.Info("Global config ", zap.Any("config values", GlobalConfig))
}

func loadConfigFromFile(file string, target interface{}) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, target); err != nil {
		return err
	}
	return nil
}

type Config struct {
	CheckItems CheckItems `yaml:"checkItems,omitempty"`
}

type CheckItems struct {
	ForbiddenPrivilegedContainer bool `yaml:"forbiddenPrivilegedContainer"`
}
