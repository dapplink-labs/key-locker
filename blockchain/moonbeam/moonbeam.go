package moonbeam

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/savour-labs/key-locker/blockchain"
	"github.com/savour-labs/key-locker/blockchain/fallback"
	"github.com/savour-labs/key-locker/config"
	"github.com/savour-labs/key-locker/crypto"
	"github.com/savour-labs/key-locker/db"
	"github.com/savour-labs/key-locker/model"
	"github.com/savour-labs/key-locker/proto/keylocker"
	"gorm.io/gorm"
)

const ChainName = "Moonbeam"

type KeyAdaptor struct {
	fallback.KeyAdaptor
	clients *KeyLockerClient
	conf    *config.Config
	repo    *model.Repo
}

func NewChainAdaptor(conf *config.Config) (blockchain.KeyAdaptor, error) {
	client, err := NewKeyLockerClient(conf)
	if err != nil {
		return nil, err
	}
	return &KeyAdaptor{
		clients: client,
		conf:    conf,
		repo:    model.NewRepo(db.InitDB(conf.Database)),
	}, nil
}

func (a *KeyAdaptor) bytesCombine(pBytes ...[]byte) []byte {
	length := len(pBytes)
	s := make([][]byte, length)
	for index := 0; index < length; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}

func (a *KeyAdaptor) GetSupportChain(req *keylocker.SupportChainReq) (*keylocker.SupportChainRep, error) {
	return &keylocker.SupportChainRep{
		Code: keylocker.ReturnCode_SUCCESS,
		Msg:  "get support chain success",
	}, nil
}

func (a *KeyAdaptor) GetSocialKey(ctx context.Context, req *keylocker.GetSocialKeyReq) (*keylocker.GetSocialKeyRep, error) {
	uuidByte := []byte(req.WalletUuid)
	var uuidByte32 [UuidSize]byte
	copy(uuidByte32[:], uuidByte)
	ret, err := a.clients.QuerySocialKey(uuidByte32)
	if err != nil {
		return nil, fmt.Errorf("GetSocialKey fail, req, %v, err: [%w]", req, err)
	}
	keyList := make([]*keylocker.SocialKey, 0)
	for _, vkey := range ret {
		keyList = append(keyList, &keylocker.SocialKey{
			Id:  "",
			Key: string(vkey),
		})
	}
	return &keylocker.GetSocialKeyRep{
		Code:    keylocker.ReturnCode_SUCCESS,
		Msg:     "get social key success",
		KeyList: keyList,
	}, nil
}

func (a *KeyAdaptor) SetSocialKey(ctx context.Context, req *keylocker.SetSocialKeyReq) (*keylocker.SetSocialKeyRep, error) {
	// get rsa key from db or generate new one
	pri, pub, encryptPriv := "", "", []byte("")
	dcrypted_pwd, err := crypto.AesDecrypt([]byte(req.Password), []byte(a.conf.AesKey))
	if err != nil {
		return nil, fmt.Errorf("decrypt password fail err: [%w]", err)
	}
	dcrypted_scode, err := crypto.AesDecrypt([]byte(req.SocialCode), []byte(a.conf.AesKey))
	if err != nil {
		return nil, fmt.Errorf("decrypt social code fail err: [%w]", err)
	}
	sec, err := a.repo.GetByUID(ctx, req.WalletUuid)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("repo.GetByUID fail, req, %v, err: [%w]", req, err)
		}
		// generate a new one
		pri, pub = crypto.NewRsa("", "").CreatePkcs8Keys(2048)
		// password  decrypt
		encryptPriv, err = crypto.AesEncrypt([]byte(pri), a.bytesCombine(dcrypted_pwd, dcrypted_scode))
		if err != nil {
			return nil, fmt.Errorf("crypto.EncryptByAes fail, req, %v, pri, %s, err: [%w]", req, pri, err)
		}
		if e := a.repo.DB.Create(&model.Secret{
			KeyUuid: req.WalletUuid,
			RsaPriv: string(encryptPriv),
			RsaPub:  pub,
		}).Error; e != nil {
			return nil, fmt.Errorf("DB.Create fail, req, %v, err: [%w]", req, e)
		}
	} else {
		priTmp, err := crypto.AesDecrypt([]byte(sec.RsaPriv), a.bytesCombine(dcrypted_pwd, dcrypted_scode))
		if err != nil {
			return nil, fmt.Errorf("crypto.DecryptByAes fail, req, %v, pri, %s, err: [%w]", req, sec.RsaPriv, err)
		}
		pri, pub = string(priTmp), sec.RsaPub
		encryptPriv = []byte(sec.RsaPriv)
	}
	// encrypt the key
	key, err := crypto.NewRsa(pub, pri).Encrypt([]byte(req.Key))
	if err != nil {
		return nil, fmt.Errorf("RSA.Encrypt fail, req, %v, err: [%w]", req, err)
	}

	uuidByte := []byte(req.WalletUuid)
	var uuidByte32 [UuidSize]byte
	copy(uuidByte32[:], uuidByte)

	if err := a.clients.AppendSocialKey(uuidByte32, [][]byte{key}); err != nil {
		return nil, err
	}

	if e := a.repo.DB.Create(&model.Key{
		KeySecret: req.Password,
		KeyUuid:   req.WalletUuid,
	}).Error; e != nil {
		return nil, fmt.Errorf("DB.Create fail, req, %v, err: [%w]", req, e)
	}
	return &keylocker.SetSocialKeyRep{
		Code: keylocker.ReturnCode_SUCCESS,
		Msg:  "set social key success",
		Pub:  pub,
		Priv: string(encryptPriv),
	}, nil
}
