package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palemoky/lucky-day/internal/model"
)

func TestLoadPrizes(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(t *testing.T) string
		wantErr   bool
		validate  func(t *testing.T, prizes []model.Prize)
	}{
		{
			name: "加载有效的奖品配置",
			setupFunc: func(t *testing.T) string {
				dir := t.TempDir()
				configContent := `
prizes:
  - id: 1
    name: "一等奖"
    count: 1
    level: 0
    probability: 0.1
  - id: 2
    name: "二等奖"
    count: 3
    level: 1
    probability: 0.2
`
				err := os.WriteFile(filepath.Join(dir, "config.yml"), []byte(configContent), 0o644)
				require.NoError(t, err)
				return dir
			},
			wantErr: false,
			validate: func(t *testing.T, prizes []model.Prize) {
				require.Len(t, prizes, 2)

				// 验证第一个奖品
				assert.Equal(t, 1, prizes[0].ID)
				assert.Equal(t, "一等奖", prizes[0].Name)
				assert.Equal(t, 1, prizes[0].Count)
				assert.Equal(t, model.PrizeLevel(0), prizes[0].Level)
				assert.Equal(t, 0.1, prizes[0].Probability)

				// 验证第二个奖品
				assert.Equal(t, 2, prizes[1].ID)
				assert.Equal(t, "二等奖", prizes[1].Name)
				assert.Equal(t, 3, prizes[1].Count)
				assert.Equal(t, model.PrizeLevel(1), prizes[1].Level)
				assert.Equal(t, 0.2, prizes[1].Probability)
			},
		},
		{
			name: "空奖品列表",
			setupFunc: func(t *testing.T) string {
				dir := t.TempDir()
				configContent := `prizes: []`
				err := os.WriteFile(filepath.Join(dir, "config.yml"), []byte(configContent), 0o644)
				require.NoError(t, err)
				return dir
			},
			wantErr: false,
			validate: func(t *testing.T, prizes []model.Prize) {
				assert.Empty(t, prizes)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupFunc(t)

			prizes, err := LoadPrizes(dir)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if tt.validate != nil {
				tt.validate(t, prizes)
			}
		})
	}
}

func TestLoadDataSourceConfig(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(t *testing.T) string
		wantErr   bool
		validate  func(t *testing.T, cfg DataSourceConfig)
	}{
		{
			name: "Excel数据源配置",
			setupFunc: func(t *testing.T) string {
				dir := t.TempDir()
				configContent := `
datasource:
  type: excel
  excel:
    path: "examples/lottery_template.xlsx"
`
				err := os.WriteFile(filepath.Join(dir, "config.yml"), []byte(configContent), 0o644)
				require.NoError(t, err)
				return dir
			},
			wantErr: false,
			validate: func(t *testing.T, cfg DataSourceConfig) {
				assert.Equal(t, "excel", cfg.Type)
				assert.Equal(t, "examples/lottery_template.xlsx", cfg.Excel.Path)
			},
		},
		{
			name: "CSV数据源配置",
			setupFunc: func(t *testing.T) string {
				dir := t.TempDir()
				configContent := `
datasource:
  type: csv
  csv:
    path: "examples/participants.csv"
`
				err := os.WriteFile(filepath.Join(dir, "config.yml"), []byte(configContent), 0o644)
				require.NoError(t, err)
				return dir
			},
			wantErr: false,
			validate: func(t *testing.T, cfg DataSourceConfig) {
				assert.Equal(t, "csv", cfg.Type)
				assert.Equal(t, "examples/participants.csv", cfg.CSV.Path)
			},
		},
		{
			name: "数据库数据源配置",
			setupFunc: func(t *testing.T) string {
				dir := t.TempDir()
				configContent := `
datasource:
  type: db
  database:
    driver: sqlite
    dsn: "lottery.db"
`
				err := os.WriteFile(filepath.Join(dir, "config.yml"), []byte(configContent), 0o644)
				require.NoError(t, err)
				return dir
			},
			wantErr: false,
			validate: func(t *testing.T, cfg DataSourceConfig) {
				assert.Equal(t, "db", cfg.Type)
				assert.Equal(t, "sqlite", cfg.Database.Driver)
				assert.Equal(t, "lottery.db", cfg.Database.DSN)
			},
		},
		{
			name: "配置文件不存在",
			setupFunc: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.setupFunc(t)

			cfg, err := LoadDataSourceConfig(dir)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}
