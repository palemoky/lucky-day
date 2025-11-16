package datasource

import (
	"fmt"

	"github.com/palemoky/lucky-day/internal/config"
	"github.com/palemoky/lucky-day/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newDBConnection 根据配置创建数据库连接
func newDBConnection(cfg config.DatabaseConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch cfg.Driver {
	case "sqlite":
		dialector = sqlite.Open(cfg.DSN)
	case "mysql":
		dialector = mysql.Open(cfg.DSN)
	case "postgres":
		dialector = postgres.Open(cfg.DSN)
	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", cfg.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 自动迁移，如果表不存在则创建
	err = db.AutoMigrate(&model.Participant{}, &model.WinningRecord{})
	if err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %w", err)
	}

	return db, nil
}

// loadParticipantsFromDB 从数据库加载所有参与者
func loadParticipantsFromDB(cfg config.DatabaseConfig) ([]model.Participant, error) {
	db, err := newDBConnection(cfg)
	if err != nil {
		return nil, err
	}

	var participants []model.Participant
	// Preload("WinningHistory") 会自动加载关联的往年中奖记录
	result := db.Preload("WinningHistory").Find(&participants)
	if result.Error != nil {
		return nil, fmt.Errorf("查询参与者失败: %w", result.Error)
	}

	fmt.Printf("成功从数据库 [%s] 加载了 %d 名参与者。\n", cfg.Driver, result.RowsAffected)
	return participants, nil
}
