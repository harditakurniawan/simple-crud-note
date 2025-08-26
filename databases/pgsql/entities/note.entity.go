package entities

type Note struct {
	Base    `gorm:"embedded"`
	Title   string  `gorm:"not null;size:100;index:idx_user_id_title,priority:2"`
	Content *string `gorm:"type:text;default:null"`
	UserID  uint    `json:"user_id" gorm:"not null;index;index:idx_user_id_title,priority:1"`
	User    User    `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName specifies the table name for Note
// func (Note) TableName() string {
// 	return "notes"
// }
