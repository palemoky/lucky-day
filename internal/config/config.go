package config

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/palemoky/lucky-day/internal/model"
)

type DataSourceConfig struct {
	Type     string         `mapstructure:"type"`
	CSV      CSVConfig      `mapstructure:"csv"`
	Excel    ExcelConfig    `mapstructure:"excel"`
	Database DatabaseConfig `mapstructure:"database"`
}

type CSVConfig struct {
	Path string `mapstructure:"path"`
}

type ExcelConfig struct {
	Path string `mapstructure:"path"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

// LoadDataSourceConfig 加载数据源配置
func LoadDataSourceConfig(path string) (DataSourceConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return DataSourceConfig{}, fmt.Errorf("无法读取配置文件: %w", err)
	}

	var config DataSourceConfig
	if err := viper.UnmarshalKey("datasource", &config); err != nil {
		return DataSourceConfig{}, fmt.Errorf("解析 datasource 配置失败: %w", err)
	}
	return config, nil
}

// LoadPrizes
func LoadPrizes(path string) ([]model.Prize, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(path)

	// 确保 viper 已经读取过配置
	if err := viper.ReadInConfig(); err != nil {
		// 如果文件不存在，可以忽略，因为可能已经被 LoadDataSourceConfig 读取过了
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("fatal error config file: %s", err)
		}
	}

	var prizes []model.Prize
	if err := viper.UnmarshalKey("prizes", &prizes); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return prizes, nil
}
