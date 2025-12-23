package datasource

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/palemoky/lucky-day/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
