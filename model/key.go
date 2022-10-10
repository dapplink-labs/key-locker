package model

import (
	"gorm.io/gorm"
)

type Key struct {
	*gorm.Model
	ID        uint64 `gorm:"primary_key; description:ID"                    json:"id"`
	KeySecret string `gorm:"index;type:varchar(256);description:KeySecret"  json:"key_secret"` // 私钥分片加密
	KeyUuid   string `gorm:"type:varchar(256);description:KeyUuid"          json:"key_uuid"`
}
