package i18n

// translations contains all UI strings in both Chinese and English
var translations = map[Language]map[string]string{
	Chinese: {
		// Application
		"app.title":         "幸运抽奖",
		"app.exit":          "抽奖结束，感谢使用！",
		"app.error":         "程序出现错误",
		"app.press_any_key": "按任意键继续...",

		// Language Selection
		"lang.select":       "请选择语言 / Select Language",
		"lang.chinese":      "中文",
		"lang.english":      "English",
		"lang.instruction":  "使用 ↑/↓ 选择，回车确认",
		"lang.instruction2": "Use ↑/↓ to select, Enter to confirm",

		// Mode Selection
		"mode.select":      "请选择抽奖模式",
		"mode.excel":       "从 Excel 文件导入",
		"mode.qr":          "二维码签到模式",
		"mode.db":          "从数据库加载",
		"mode.instruction": "使用 ↑/↓ 选择，回车确认，q 退出",

		// Prize Selection
		"prize.title":       "请选择要抽取的奖项",
		"prize.instruction": "使用 ↑/↓ 选择，回车开始抽奖，q 退出",
		"prize.remaining":   "剩余",
		"prize.total":       "总共",
		"prize.all_drawn":   "已抽完",

		// Drawing
		"draw.title":       "正在抽奖...",
		"draw.instruction": "按任意键停止",
		"draw.rolling":     "滚动中",

		// Winners
		"winner.title":        "恭喜中奖！",
		"winner.instruction":  "按任意键返回奖项选择",
		"winner.list_title":   "中奖名单",
		"winner.no_winners":   "暂无中奖者",
		"winner.prize":        "奖项",
		"winner.name":         "姓名",
		"winner.save_success": "中奖名单已保存到 Excel",
		"winner.save_failed":  "保存中奖名单失败",

		// QR Check-in
		"qr.title":              "二维码签到",
		"qr.instruction":        "请用手机扫描二维码进行签到",
		"qr.url":                "签到地址",
		"qr.qr_file":            "二维码文件",
		"qr.count":              "已签到人数",
		"qr.start":              "按 Enter 开始抽奖",
		"qr.save":               "按 s 保存签到名单",
		"qr.quit":               "按 q 退出",
		"qr.saved":              "签到名单已保存",
		"qr.checkin_success":    "签到成功！",
		"qr.checkin_failed":     "签到失败",
		"qr.name_required":      "请输入姓名",
		"qr.name_placeholder":   "请输入您的姓名",
		"qr.dept_placeholder":   "部门（可选）",
		"qr.submit":             "提交签到",
		"qr.ready":              "二维码签到已准备就绪！",
		"qr.hint":               "请在其他设备上打开二维码图片，让参与者扫码签到",
		"qr.press_enter":        "签到完成后按 Enter 键开始抽奖",
		"qr.no_participants":    "没有人签到",
		"qr.total_participants": "共有签到人数",

		// Data Source
		"data.loading":         "正在加载数据...",
		"data.load_success":    "成功加载",
		"data.load_failed":     "加载失败",
		"data.participants":    "名参与者",
		"data.prizes":          "个奖项",
		"data.source_csv":      "数据源: CSV 文件",
		"data.source_excel":    "数据源: Excel 文件",
		"data.source_db":       "数据源: 数据库",
		"data.empty_list":      "参与者名单为空",
		"data.config_error":    "配置文件错误",
		"data.excel_not_found": "Excel 文件不存在",

		// Errors
		"error.unknown":        "未知错误",
		"error.file_not_found": "文件不存在",
		"error.invalid_config": "配置无效",
		"error.network":        "网络错误",
		"error.server_start":   "服务器启动失败",

		// Footer
		"footer.help": "帮助: ↑/↓ 选择 | Enter 确认 | q 退出",
	},

	English: {
		// Application
		"app.title":         "Lucky Draw",
		"app.exit":          "Lottery ended, thank you!",
		"app.error":         "An error occurred",
		"app.press_any_key": "Press any key to continue...",

		// Language Selection
		"lang.select":       "请选择语言 / Select Language",
		"lang.chinese":      "中文",
		"lang.english":      "English",
		"lang.instruction":  "使用 ↑/↓ 选择，回车确认",
		"lang.instruction2": "Use ↑/↓ to select, Enter to confirm",

		// Mode Selection
		"mode.select":      "Select Lottery Mode",
		"mode.excel":       "Load from Excel File",
		"mode.qr":          "QR Code Check-in Mode",
		"mode.db":          "Load from Database",
		"mode.instruction": "Use ↑/↓ to select, Enter to confirm, q to quit",

		// Prize Selection
		"prize.title":       "Select a Prize to Draw",
		"prize.instruction": "Use ↑/↓ to select, Enter to start, q to quit",
		"prize.remaining":   "Remaining",
		"prize.total":       "Total",
		"prize.all_drawn":   "All Drawn",

		// Drawing
		"draw.title":       "Drawing...",
		"draw.instruction": "Press any key to stop",
		"draw.rolling":     "Rolling",

		// Winners
		"winner.title":        "Congratulations!",
		"winner.instruction":  "Press any key to return",
		"winner.list_title":   "Winners List",
		"winner.no_winners":   "No winners yet",
		"winner.prize":        "Prize",
		"winner.name":         "Name",
		"winner.save_success": "Winners saved to Excel",
		"winner.save_failed":  "Failed to save winners",

		// QR Check-in
		"qr.title":              "QR Code Check-in",
		"qr.instruction":        "Scan QR code with your phone to check in",
		"qr.url":                "Check-in URL",
		"qr.qr_file":            "QR Code File",
		"qr.count":              "Checked-in Count",
		"qr.start":              "Press Enter to start lottery",
		"qr.save":               "Press s to save check-in list",
		"qr.quit":               "Press q to quit",
		"qr.saved":              "Check-in list saved",
		"qr.checkin_success":    "Check-in successful!",
		"qr.checkin_failed":     "Check-in failed",
		"qr.name_required":      "Name is required",
		"qr.name_placeholder":   "Enter your name",
		"qr.dept_placeholder":   "Department (optional)",
		"qr.submit":             "Submit Check-in",
		"qr.ready":              "QR Code check-in is ready!",
		"qr.hint":               "Open the QR code image on another device for participants to scan",
		"qr.press_enter":        "Press Enter to start lottery after check-in is complete",
		"qr.no_participants":    "No participants checked in",
		"qr.total_participants": "Total participants",

		// Data Source
		"data.loading":         "Loading data...",
		"data.load_success":    "Successfully loaded",
		"data.load_failed":     "Failed to load",
		"data.participants":    "participants",
		"data.prizes":          "prizes",
		"data.source_csv":      "Data source: CSV file",
		"data.source_excel":    "Data source: Excel file",
		"data.source_db":       "Data source: Database",
		"data.empty_list":      "Participant list is empty",
		"data.config_error":    "Configuration error",
		"data.excel_not_found": "Excel file not found",

		// Errors
		"error.unknown":        "Unknown error",
		"error.file_not_found": "File not found",
		"error.invalid_config": "Invalid configuration",
		"error.network":        "Network error",
		"error.server_start":   "Failed to start server",

		// Footer
		"footer.help": "Help: ↑/↓ Select | Enter Confirm | q Quit",
	},
}
