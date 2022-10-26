package model

type SnakeLadder struct {
	Current      []int   `json:"current"`        // index:0自己位置, 其他机器人队友位置;value -1:未开始
	Map          []int   `json:"map"`            // 地图 ;0空 1蛇头 2蛇尾 3梯子头 4梯子尾
	SavePlayBack [][]int `json:"save_play_back"` // 保存每次的骰子结果，用于回放
	WinPlayer    int     `json:"win_player"`     // 获胜人编号 0未分出胜负 1自己 >1机器人
}

type RollSnakeLadder struct {
	Points      []int        `json:"points"`       // 玩家+机器人骰子的结果
	SnakeLadder *SnakeLadder `json:"snake_ladder"` // 位置 地图等信息
}
