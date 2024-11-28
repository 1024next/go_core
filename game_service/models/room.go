package models

// Room 模型，表示一个房间
type Room struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Players  string `json:"players"`   // 玩家列表，以JSON字符串形式保存
	MaxSeats int    `json:"max_seats"` // 最大座位数
	GameID   uint   `json:"game_id"`   // 关联的游戏ID
}

// Room 表示房间数据的数据库模型
func (Room) TableName() string {
	return "rooms" // 表名
}
