package models

// Game 模型，表示一个游戏
type Game struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	RoomCount int    `json:"room_count"`
}

// Game 表示游戏数据的数据库模型
func (Game) TableName() string {
	return "games" // 表名
}
