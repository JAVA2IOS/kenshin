package models

// 导航栏结构体
type NaviAction struct {
	Action string // 动作
	Name   string // 名字
}

// 导航栏结构
func NaviActions() []NaviAction {

	return []NaviAction{{Action: "jd_", Name: "京东"}, {Action: "pdd_", Name: "拼多多"}}
}
