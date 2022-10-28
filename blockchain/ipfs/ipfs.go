package ipfs

import (
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

const ChainName = "Ipfs"

type KeyAdaptor struct {
	fallback.KeyAdaptor
	repo       *model.Repo
	ipfsClient *Client
}

func NewChainAdaptor(conf *config.Config) (blockchain.KeyAdaptor, error) {
	ipfsClient, err := New(context.Background(), conf.Fullnode.Ipfs.NetworkNode, conf.Fullnode.Ipfs.RepoPath)
	if err != nil {
		return nil, err
	}
	return &KeyAdaptor{
		repo:       model.NewRepo(db.InitDB(conf.Database)),
		ipfsClient: ipfsClient,
	}, nil
}

func (a *KeyAdaptor) GetSupportChain(req *keylocker.SupportChainReq) (*keylocker.SupportChainRep, error) {
	return &keylocker.SupportChainRep{
		Code: keylocker.ReturnCode_SUCCESS,
		Msg:  "get support chain success",
	}, nil
}

// GetSocialKey
// 1. req.uuid 取到链上的存储 req.key 的文件，
// 2. 解密返回就行
func (a *KeyAdaptor) GetSocialKey(ctx context.Context, req *keylocker.GetSocialKeyReq) (*keylocker.GetSocialKeyRep, error) {
	// get file from ipfs
	ret, err := a.ipfsClient.GetFile(ctx, req.FileCid)
	if err != nil {
		return nil, fmt.Errorf("ipfsClient.GetFile fail, req, %v, err: [%w]", req, err)
	}

	// get rsa key from db
	sec, err := a.repo.GetByUID(ctx, req.Uuid)
	if err != nil {
		return nil, fmt.Errorf("repo.GetByUID fail, req, %v, err: [%w]", req, err)
	}

	// Decrypt rsa key
	pri, err := crypto.DecryptByAes(sec.RsaPriv)
	if err != nil {
		return nil, fmt.Errorf("crypto.DecryptByAes fail, req, %v, err: [%w]", req, err)
	}

	// Decrypt the content from ipfs
	key, err := crypto.NewRsa(sec.RsaPub, string(pri)).Decrypt(ret)
	if err != nil {
		return nil, fmt.Errorf("RSA.Decrypt fail, req, %v, err: [%w]", req, err)
	}

	return &keylocker.GetSocialKeyRep{
		Code: keylocker.ReturnCode_SUCCESS,
		Msg:  "get ipfs social key success",
		KeyList: []*keylocker.SocialKey{&keylocker.SocialKey{
			Id:  "",
			Key: string(key),
		}},
	}, nil
}

// SetSocialKey
// 1. 如果对应 uuid(req.uuid) 没有 rsa 密钥对，生成 rsa 密钥对，生成 RSA 密钥对处理，私钥，用用户密码(req.password)进行 AES 加密存储， 公钥明文存储, 如果有直接使用
// 2. 用 rsa 私钥对 key(req.key 是用户上传的一个私钥) 加密，加密 key 调用 ipfs 上传
// 3. 返回加密的 RSA 的私钥和明文的 RSA 公钥匙, 加密方式，IPFS 对应的 CID
func (a *KeyAdaptor) SetSocialKey(ctx context.Context, req *keylocker.SetSocialKeyReq) (*keylocker.SetSocialKeyRep, error) {
	// get rsa key from db or generate new one
	pri, pub := "", ""
	sec, err := a.repo.GetByUID(ctx, req.Uuid)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("repo.GetByUID fail, req, %v, err: [%w]", req, err)
		}

		// generate a new one
		pri, pub = crypto.NewRsa("", "").CreatePkcs8Keys(2048)
		priByAES, err := crypto.EncryptByAes([]byte(pri))
		if err != nil {
			return nil, fmt.Errorf("crypto.EncryptByAes fail, req, %v, pri, %s, err: [%w]", req, pri, err)
		}
		if e := a.repo.DB.Create(&model.Secret{
			KeyUuid: req.Uuid,
			RsaPriv: priByAES,
			RsaPub:  pub,
		}).Error; e != nil {
			return nil, fmt.Errorf("DB.Create fail, req, %v, err: [%w]", req, e)
		}
	} else {
		priTmp, err := crypto.DecryptByAes(sec.RsaPriv)
		if err != nil {
			return nil, fmt.Errorf("crypto.DecryptByAes fail, req, %v, pri, %s, err: [%w]", req, sec.RsaPriv, err)
		}

		pri, pub = string(priTmp), sec.RsaPub
	}

	// encrypt the key
	key, err := crypto.NewRsa(pub, pri).Encrypt([]byte(req.Key))
	if err != nil {
		return nil, fmt.Errorf("RSA.Encrypt fail, req, %v, err: [%w]", req, err)
	}
	cid, err := a.ipfsClient.AddFile(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("ipfsClient.AddFile fail, req, %v, err: [%w]", req, err)
	}

	// insert into db
	if e := a.repo.DB.Create(&model.Key{
		KeySecret: req.Password,
		KeyCID:    cid,
		KeyUuid:   req.Uuid,
	}).Error; e != nil {
		return nil, fmt.Errorf("DB.Create fail, req, %v, err: [%w]", req, e)
	}

	return &keylocker.SetSocialKeyRep{
		Code:    keylocker.ReturnCode_SUCCESS,
		Msg:     "set ipfs social key success",
		Pub:     pub,
		Priv:    pri,
		FileCid: cid,
	}, nil
}
