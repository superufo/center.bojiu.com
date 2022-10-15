package storage

//代理类型
const (
	AGENT_LEVEL = 1 + iota
	AGENT_PYRAMID
)

//控制状态
const (
	CTRL_NORMAL = iota
	CTRL_KILL
	CTRL_GIVE
)

type AgentLevel struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
	Ratio string `json:"ratio"`
}
