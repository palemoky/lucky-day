# Lucky Day

![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Test Status](https://img.shields.io/badge/tests-passing-brightgreen.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/palemoky/lucky-day)](https://goreportcard.com/report/github.com/palemoky/lucky-day)

一款功能强大、高度可配置的终端抽奖应用程序。专为公司年会、线上活动、社区抽奖等场景设计，提供公平、透明且富有观赏性的抽奖体验。

![Lucky Day Screenshot](docs/screenshot.png)

---

## ✨ 功能特性

- **美观的终端界面 (TUI)**: 基于 [Bubble Tea](https://github.com/charmbracelet/bubbletea) 和 [Lipgloss](https://github.com/charmbracelet/lipgloss) 构建，提供流畅的交互体验和动态抽奖动画。
- **高度可配置的奖品**: 所有奖品信息（名称、数量、等级）均在 `config.yml` 文件中定义，无需修改代码即可调整。
- **灵活的数据源**:
  - 支持从 **CSV 文件** 快速导入参与者名单。
  - 通过 **GORM** 集成，原生支持 **SQLite**, **MySQL**, **PostgreSQL** 等多种数据库，轻松应对海量用户数据。
- **智能的加权抽奖算法**:
  - 可根据参与者的**往年中奖历史**动态调整中奖权重，让新人有更多机会。
  - 权重算法考虑了中奖年份和奖品等级，实现更精细的公平性控制。
- **实时的抽奖动画**: 抽取多人奖项时，屏幕上会同时滚动多个名字，按下任意键即可“定格”中奖者，增加现场紧张感和趣味性。
- **持久化的中奖名单**: 右侧面板会实时展示所有已中奖人员名单，一目了然。
- **健壮且经过测试**:
  - 核心逻辑拥有完整的**单元测试** (基于 `testify`)。
  - 通过 **Go Fuzzing (模糊测试)** 探索边缘情况，确保抽奖引擎在各种极端输入下依然稳定。
- **跨平台**: 基于 Go 语言开发，可轻松编译并在 Windows, macOS, Linux 等主流操作系统上运行。
