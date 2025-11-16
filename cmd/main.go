package main

import (
	"fmt"
	"log"
	"os"

	"github.com/palemoky/lucky-day/internal/config"
	"github.com/palemoky/lucky-day/internal/datasource"
	"github.com/palemoky/lucky-day/internal/lottery"
	"github.com/palemoky/lucky-day/internal/tui"
)

func main() {
	// 1. 加载配置
	prizes, err := config.LoadPrizes(".")
	if err != nil {
		log.Fatalf("加载奖品配置失败: %v", err)
	}

	dsCfg, err := config.LoadDataSourceConfig(".")
	if err != nil {
		log.Fatalf("加载数据源配置失败: %v", err)
	}

	// 2. 根据配置加载数据
	participants, err := datasource.LoadParticipants(dsCfg)
	if err != nil {
		log.Fatalf("加载参与者名单失败: %v", err)
	}
	if len(participants) == 0 {
		log.Fatal("参与者名单为空，程序无法运行。")
	}

	// 3. 初始化抽奖引擎
	engine := lottery.NewEngine(participants, prizes)

	// 4. 启动并运行TUI
	if err := tui.StartTUI(engine); err != nil {
		fmt.Printf("程序出现错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("抽奖结束，感谢使用！")
}