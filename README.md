# 乐城堡社招笔试题目-蛇梯游戏
游戏的参与者是1玩家和N机器人 ，现在设置的机器人数量是4



## 服务简介

利用自定义格式的TCP协议跟前端交互，socket断开连接时候保存数据到MongoDB

配置一下config.json就能运行服务端

客户端在internal/client.go

日志 log/2022-10-26.log

## API

**通信协议**

Tcp

 

**请求方法**

Tcpx协议

 

**字符编码**

UTF-8 

## 公共参数

无

## 鉴权机制

无

## 响应结果

**成功响应**

当 code 状态码为 ``00000`` 时，表明调用成功。

### 协议说明

### 协议号

| 协议 | **说明**     |
| ---- | ------------ |
| 1    | 登陆         |
| 2    | 进行蛇梯游戏 |

### 协议字段

协议 1 略



协议 2 进行蛇梯游戏

- Req

  ```go
  Action int `json:"action"` // 1:初始化棋盘 初始化N个队友 2:掷色子 3：回放
  ```

  ```json
  {"action":1}
  {"action":2}
  {"action":3}
  ```

  

- Resp

  ```go
  Code     string `json:"code,omitempty"` // 状态码 成功 "00000" ； 游戏结束 "20001"
  Result   *model.RollSnakeLadder `json:"result,omitempty"` // 结果
  Playback [][]int                `json:"playback,omitempty"` // 回放记录
  ```

  ```go
  type RollSnakeLadder struct {
  	Points      []int        `json:"points"`       // 玩家+机器人骰子的结果
  	SnakeLadder *SnakeLadder `json:"snake_ladder"` // 位置 地图等信息
  }
  ```

  初始化棋盘

  ```json
  {"code":"00000"}
  ```

  掷色子

  ```json
  {"code":"00000","result":{"points":[4,2,4,2,2],"snake_ladder":{"current":[3,1,3,1,1],"map":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"save_play_back":[[4],[2],[4],[2],[2]],"win_player":0}}}
  ```

  游戏结束

  ``` json
  {"code":"00000","result":{"points":[2,4,3,6,5],"snake_ladder":{"current":[98,94,97,97,99],"map":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"save_play_back":[[4,6,6,6,3,4,3,2,5,3,5,6,2,1,3,6,5,5,1,5,1,4,5,2,6,4,2],[2,5,1,1,4,3,2,3,5,6,2,6,4,1,1,4,6,2,6,4,4,5,1,4,3,6,4],[4,3,6,3,4,5,3,2,4,2,6,6,6,3,4,2,3,3,4,6,4,1,2,1,6,2,3],[2,1,4,5,3,3,3,3,1,1,6,3,5,4,3,5,1,4,4,6,2,5,5,5,4,4,6],[2,6,4,3,4,1,6,5,5,3,1,6,6,1,6,1,4,3,3,1,4,3,5,5,5,2]],"win_player":5}}}
  ```

  回放

  ```json
  {"code":"00000","playback":[[4,6,6,6,3,4,3,2,5,3,5,6,2,1,3,6,5,5,1,5,1,4,5,2,6,4,2],[2,5,1,1,4,3,2,3,5,6,2,6,4,1,1,4,6,2,6,4,4,5,1,4,3,6,4],[4,3,6,3,4,5,3,2,4,2,6,6,6,3,4,2,3,3,4,6,4,1,2,1,6,2,3],[2,1,4,5,3,3,3,3,1,1,6,3,5,4,3,5,1,4,4,6,2,5,5,5,4,4,6],[2,6,4,3,4,1,6,5,5,3,1,6,6,1,6,1,4,3,3,1,4,3,5,5,5,2]]}
  ```

  