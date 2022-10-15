package model

import "gorm.io/gorm"

type Secret struct {
	*gorm.Model
	ID      uint64 `gorm:"primary_key; description:ID"                    json:"id"`
	KeyUuid string `gorm:"type:varchar(256);description:KeyUuid"          json:"key_uuid"`
	RsaPriv string `gorm:"type:varchar(2048);description:RsaPriv"         json:"rsa_priv"`
	RsaPub  string `gorm:"type:varchar(2048);description:RsaPub"          json:"rsa_pub"`
}
