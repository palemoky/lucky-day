# Lucky Day 🎉

一款功能强大、高度可配置的终端抽奖应用程序。专为公司年会、线上活动、社区抽奖等场景设计，提供公平、透明且富有观赏性的抽奖体验。

![Lucky Day Screenshot](docs/screenshot.png)

---

## ✨ 核心特性

### 🎨 美观的用户界面

- **全屏 TUI 界面**：基于 [Bubble Tea](https://github.com/charmbracelet/bubbletea) 和 [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **流畅的交互体验**：无闪烁的界面切换
- **实时抽奖动画**：多人奖项同时滚动名字，按键定格中奖者
- **中奖名单展示**：右侧面板实时显示所有中奖人员

### 🌍 国际化支持

- **中英双语**：启动时选择语言
- **完整翻译**：所有 UI 文本支持中英文
- **灵活切换**：可随时更改语言设置

### 📊 多种数据源

- **Excel 文件** (推荐)：
  - 三个 Sheet 管理：奖品设置、参与者、中奖历史
  - 支持 10 万+ 参与者
  - 自动保存中奖记录
- **二维码签到**：
  - 实时移动端签到
  - 自动生成 QR 码
  - HTTP 服务器接收签到
- **数据库**：
  - 支持 SQLite、MySQL、PostgreSQL
  - 通过 GORM 集成

### 🎯 智能抽奖算法

- **加权抽奖**：根据往年中奖历史动态调整权重
- **公平性保证**：新人有更多中奖机会
- **精细控制**：考虑中奖年份和奖品等级

### 🔧 高度可配置

- **奖品配置**：`config.yml` 定义所有奖品信息
- **灵活数据源**：支持 CSV、Excel、数据库
- **无需编码**：所有配置通过配置文件完成

### 🛡️ 健壮可靠

- **完整测试**：核心逻辑有单元测试
- **模糊测试**：Go Fuzzing 探索边缘情况
- **代码质量**：通过 golangci-lint 检查
- **Pre-commit Hooks**：自动代码格式化和检查

### 🌐 跨平台

- 支持 Windows、macOS、Linux
- 单一二进制文件，无需依赖

---

## 🚀 快速开始

### 安装

#### 方式一：下载预编译二进制文件（推荐）

从 [GitHub Releases](https://github.com/palemoky/lucky-day/releases) 下载对应平台的压缩包。

#### 方式二：从源码构建

```bash
# 克隆项目
git clone --depth 1 https://github.com/palemoky/lucky-day.git
cd lucky-day

# 构建
go build -o lottery ./cmd
```

### 快速开始

1. **下载并解压**压缩包到任意目录

2. **编辑参与者名单**：

   - 打开 `examples/lottery_template.xlsx`
   - 在 "Participants" Sheet 中添加参与者信息
   - 在 "Prizes" Sheet 中配置奖品（可选）

3. **运行抽奖**：

   ```bash
   # Linux/macOS
   ./lottery

   # Windows
   lottery.exe
   ```

4. **选择语言** → **选择模式** → **开始抽奖**！

5. **操作流程**：
   - 选择语言（中文/English）
   - 选择模式（Excel/二维码/数据库）
   - 选择奖项
   - 按 Enter 开始抽奖
   - 按任意键停止

---

## 📖 使用指南

### Excel 模式

**适用场景**：公司年会、固定参与者名单

**Excel 文件结构**：

| Sheet        | 说明       | 列                                                         |
| ------------ | ---------- | ---------------------------------------------------------- |
| Prizes       | 奖品配置   | ID, Name(CN), Name(EN), Count, Level, Probability          |
| Participants | 参与者名单 | ID, Name, Department, Email                                |
| Winners      | 中奖历史   | Draw Time, Prize Name, Winner ID, Winner Name, Prize Level |

**配置**：

```yaml
datasource:
  type: excel
  excel:
    path: "examples/lottery_template.xlsx"
```

### 二维码签到模式

**适用场景**：现场活动、临时参与者

**使用步骤**：

1. 选择"二维码签到模式"
2. 程序生成 `checkin_qr.png`
3. 在投影仪/大屏幕上展示二维码
4. 参与者扫码签到
5. 按 Enter 开始抽奖

**特点**：

- 实时签到统计
- 移动端友好界面
- 自动分配参与者 ID
- 签到完成后服务器自动关闭

### 数据库模式

**适用场景**：大型活动、与现有系统集成

**支持的数据库**：

- SQLite（本地文件）
- MySQL
- PostgreSQL

**配置示例**：

```yaml
datasource:
  type: db
  database:
    driver: sqlite
    dsn: "lottery.db"
```

---

## 🎁 奖品配置

在 `config.yml` 中配置奖品：

```yaml
prizes:
  - id: 1
    name_cn: "特等奖"
    name_en: "Grand Prize"
    count: 1
    level: 0
    probability: 0.01

  - id: 2
    name_cn: "一等奖"
    name_en: "First Prize"
    count: 3
    level: 1
    probability: 0.05
```

**字段说明**：

- `id`: 奖品唯一标识
- `name_cn/name_en`: 中英文名称
- `count`: 奖品数量
- `level`: 奖品等级（0=特等奖，数字越大等级越低）
- `probability`: 中奖概率（0.0-1.0）

---

## 🔧 高级配置

### 数据源切换

```yaml
datasource:
  type: excel # 可选: excel, db, csv

  excel:
    path: "examples/lottery_template.xlsx"

  database:
    driver: sqlite
    dsn: "lottery.db"

  csv:
    path: "participants.csv"
```

### 国际化

程序启动时选择语言，所有 UI 文本自动切换。

---

## 🛠️ 开发

### 环境要求

- Go 1.24+
- golangci-lint（可选，用于代码检查）
- pre-commit（可选，用于 Git hooks）

### 安装开发工具

```bash
# 安装 pre-commit hooks
pre-commit install
pre-commit install --hook-type pre-push

# 运行代码检查
golangci-lint run

# 运行测试
go test ./...
```

### 代码规范

项目使用以下工具确保代码质量：

- `golangci-lint`: 代码静态分析
- `gofumpt`: 代码格式化
- `goimports`: Import 排序
- Pre-commit hooks: 自动检查和格式化
