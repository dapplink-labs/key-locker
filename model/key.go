package model

import (
	"gorm.io/gorm"
)

type Key struct {
	KeySecret string `gorm:"type:text;description:KeySecret; comment: uid的rsa私钥"    json:"key_secret"`
	KeyCID    string `gorm:"type:varchar(256);;description:KeyCID; comment: key对应的ipfs CID"    json:"key_cid"`
	KeyUuid   string `gorm:"index;type:varchar(256);description:KeyUuid; comment: 用户ID"    json:"key_uuid"`
	*gorm.Model
}

type Repo struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{DB: db}
}
