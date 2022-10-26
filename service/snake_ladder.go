package service

import (
	"fmt"
	"mania/model"
	"mania/util/rand"
)

const MAP_LENGTH = 100
const ROBOT_NUMBER = 4
const WIN_INDEX = 99

// 初始化地图
func InitSnakeLadderMap(usr *model.User) {
	usr.SnakeLadder = &model.SnakeLadder{
		Map:          make([]int, MAP_LENGTH),
		Current:      make([]int, ROBOT_NUMBER+1),
		SavePlayBack: make([][]int, ROBOT_NUMBER+1),
		WinPlayer:    0,
	}

	for i := range usr.SnakeLadder.Current {
		usr.SnakeLadder.Current[i] = -1
	}
}

func Roll(usr *model.User) (r *model.RollSnakeLadder, err error) {
	if usr.SnakeLadder == nil {
		err = fmt.Errorf("game not initialized")
		return
	}

	if usr.SnakeLadder.WinPlayer > 0 {
		err = fmt.Errorf("the game is over")
		return
	}

	r = new(model.RollSnakeLadder)
	r.Points = make([]int, ROBOT_NUMBER+1)

	for i := range usr.SnakeLadder.Current {
		rollPoint := rand.RandomInt(6) + 1
		r.Points[i] = rollPoint

		usr.SnakeLadder.Current[i] += rollPoint
		current := usr.SnakeLadder.Current[i]
		// 游戏结束
		if current == WIN_INDEX {
			usr.SnakeLadder.WinPlayer = i + 1
			r.SnakeLadder = usr.SnakeLadder
			return
			// 回退
		} else if current > WIN_INDEX {
			usr.SnakeLadder.Current[i] = WIN_INDEX - (current - WIN_INDEX)
		}

		usr.SnakeLadder.SavePlayBack[i] = append(usr.SnakeLadder.SavePlayBack[i], rollPoint)
	}

	r.SnakeLadder = usr.SnakeLadder
	return
}
