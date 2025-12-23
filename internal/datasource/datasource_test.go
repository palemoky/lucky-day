package datasource

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palemoky/lucky-day/internal/config"
	"github.com/palemoky/lucky-day/internal/model"
)

func TestLoadPrizesFromExcel(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(t *testing.T) string
		wantErr   bool
		validate  func(t *testing.T, prizes []model.Prize)
	}{
		{
			name: "加载有效的Excel奖品数据",
			setupFunc: func(t *testing.T) string {
				// 使用项目中的示例文件
				return "../../examples/lottery_template.xlsx"
			},
			wantErr: false,
			validate: func(t *testing.T, prizes []model.Prize) {
				// 验证至少有一些奖品
				assert.NotEmpty(t, prizes)
				// 验证奖品结构
				for _, prize := range prizes {
					assert.NotZero(t, prize.ID)
					assert.NotEmpty(t, prize.Name)
					assert.NotZero(t, prize.Count)
				}
			},
		},
		{
			name: "文件不存在",
			setupFunc: func(t *testing.T) string {
				return "nonexistent.xlsx"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setupFunc(t)

			prizes, err := LoadPrizesFromExcel(path)

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

func TestLoadParticipantsFromExcel(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(t *testing.T) string
		wantErr   bool
		validate  func(t *testing.T, participants []model.Participant)
	}{
		{
			name: "加载有效的Excel参与者数据",
			setupFunc: func(t *testing.T) string {
				return "../../examples/lottery_template.xlsx"
			},
			wantErr: false,
			validate: func(t *testing.T, participants []model.Participant) {
				// 验证至少有一些参与者
				assert.NotEmpty(t, participants)
				// 验证参与者结构
				for _, p := range participants {
					assert.NotZero(t, p.ID)
					assert.NotEmpty(t, p.Name)
				}
			},
		},
		{
			name: "文件不存在",
			setupFunc: func(t *testing.T) string {
				return "nonexistent.xlsx"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setupFunc(t)

			participants, err := LoadParticipantsFromExcel(path)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if tt.validate != nil {
				tt.validate(t, participants)
			}
		})
	}
}

func TestSaveWinnersToExcel(t *testing.T) {
	tests := []struct {
		name    string
		winners []Winner
		wantErr bool
	}{
		{
			name: "保存中奖者到Excel",
			winners: []Winner{
				{
					DrawTime:   time.Now(),
					PrizeName:  "一等奖",
					WinnerID:   1,
					WinnerName: "张三",
					PrizeLevel: 0,
				},
				{
					DrawTime:   time.Now(),
					PrizeName:  "二等奖",
					WinnerID:   2,
					WinnerName: "李四",
					PrizeLevel: 1,
				},
			},
			wantErr: false,
		},
		{
			name:    "空中奖者列表",
			winners: []Winner{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用临时文件
			tmpDir := t.TempDir()
			outputPath := filepath.Join(tmpDir, "winners.xlsx")

			// 先创建一个模板文件
			err := CreateExcelTemplate(outputPath)
			require.NoError(t, err)

			// 保存中奖者
			err = SaveWinnersToExcel(outputPath, tt.winners)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// 验证文件仍然存在
			_, err = os.Stat(outputPath)
			assert.NoError(t, err, "输出文件应该存在")
		})
	}
}

func TestCreateExcelTemplate(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "创建Excel模板",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用临时文件
			tmpDir := t.TempDir()
			outputPath := filepath.Join(tmpDir, "template.xlsx")

			err := CreateExcelTemplate(outputPath)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// 验证文件已创建
			_, err = os.Stat(outputPath)
			assert.NoError(t, err, "模板文件应该存在")

			// 尝试从创建的模板中加载数据
			prizes, err := LoadPrizesFromExcel(outputPath)
			require.NoError(t, err)
			assert.NotEmpty(t, prizes, "模板应该包含示例奖品")

			participants, err := LoadParticipantsFromExcel(outputPath)
			require.NoError(t, err)
			assert.NotEmpty(t, participants, "模板应该包含示例参与者")
		})
	}
}

func TestLoadParticipantsFromCSV(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(t *testing.T) string
		wantErr   bool
		validate  func(t *testing.T, participants []model.Participant)
	}{
		{
			name: "加载有效的CSV文件",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				csvPath := filepath.Join(tmpDir, "participants.csv")

				content := `ID,Name,Department,Email
1,张三,技术部,zhangsan@example.com
2,李四,市场部,lisi@example.com
3,王五,人事部,wangwu@example.com
`
				err := os.WriteFile(csvPath, []byte(content), 0o644)
				require.NoError(t, err)
				return csvPath
			},
			wantErr: false,
			validate: func(t *testing.T, participants []model.Participant) {
				assert.Len(t, participants, 3)
				assert.Equal(t, 1, participants[0].ID)
				assert.Equal(t, "张三", participants[0].Name)
				assert.Equal(t, 2, participants[1].ID)
				assert.Equal(t, "李四", participants[1].Name)
			},
		},
		{
			name: "空CSV文件",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				csvPath := filepath.Join(tmpDir, "empty.csv")
				content := `ID,Name
`
				err := os.WriteFile(csvPath, []byte(content), 0o644)
				require.NoError(t, err)
				return csvPath
			},
			wantErr: false,
			validate: func(t *testing.T, participants []model.Participant) {
				assert.Empty(t, participants)
			},
		},
		{
			name: "文件不存在",
			setupFunc: func(t *testing.T) string {
				return "nonexistent.csv"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setupFunc(t)

			participants, err := loadParticipantsFromCSV(path)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if tt.validate != nil {
				tt.validate(t, participants)
			}
		})
	}
}

func TestLoadParticipants(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(t *testing.T) config.DataSourceConfig
		wantErr   bool
		validate  func(t *testing.T, participants []model.Participant)
	}{
		{
			name: "从Excel加载",
			setupFunc: func(t *testing.T) config.DataSourceConfig {
				return config.DataSourceConfig{
					Type: "excel",
					Excel: config.ExcelConfig{
						Path: "../../examples/lottery_template.xlsx",
					},
				}
			},
			wantErr: false,
			validate: func(t *testing.T, participants []model.Participant) {
				assert.NotEmpty(t, participants)
			},
		},
		{
			name: "从CSV加载",
			setupFunc: func(t *testing.T) config.DataSourceConfig {
				tmpDir := t.TempDir()
				csvPath := filepath.Join(tmpDir, "test.csv")
				content := `ID,Name
1,测试用户
`
				err := os.WriteFile(csvPath, []byte(content), 0o644)
				require.NoError(t, err)

				return config.DataSourceConfig{
					Type: "csv",
					CSV: config.CSVConfig{
						Path: csvPath,
					},
				}
			},
			wantErr: false,
			validate: func(t *testing.T, participants []model.Participant) {
				assert.Len(t, participants, 1)
				assert.Equal(t, "测试用户", participants[0].Name)
			},
		},
		{
			name: "未知数据源类型",
			setupFunc: func(t *testing.T) config.DataSourceConfig {
				return config.DataSourceConfig{
					Type: "unknown",
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.setupFunc(t)

			participants, err := LoadParticipants(cfg)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if tt.validate != nil {
				tt.validate(t, participants)
			}
		})
	}
}
