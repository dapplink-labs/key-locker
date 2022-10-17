package model

import (
	"context"

	"gorm.io/gorm"
)

type Secret struct {
	*gorm.Model
	KeyUuid string `gorm:"type:varchar(256);description:KeyUuid;comment:用户ID"          json:"key_uuid"`
	RsaPriv string `gorm:"type:text;description:RsaPriv;comment:RSA私钥"         json:"rsa_priv"`
	RsaPub  string `gorm:"type:text;description:RsaPub;comment:RSA公钥"          json:"rsa_pub"`
}

func (r *Repo) GetByUID(ctx context.Context, uid string) (*Secret, error) {
	res := new(Secret)
	if err := r.DB.WithContext(ctx).Where("key_uuid = ?", uid).First(res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
