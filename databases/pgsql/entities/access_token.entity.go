package entities

type AccessToken struct {
	Base   `gorm:"embedded"`
	UserID uint   `json:"user_id" gorm:"not null;index"`
	User   User   `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Token  string `json:"token" gorm:"type:text;not null;uniqueIndex"`
}
