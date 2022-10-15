package ipfs

import (
	"github.com/savour-labs/key-locker/blockchain"
	"github.com/savour-labs/key-locker/blockchain/fallback"
	"github.com/savour-labs/key-locker/blockchain/multiclient"
	"github.com/savour-labs/key-locker/config"
	"github.com/savour-labs/key-locker/proto/common"
	"github.com/savour-labs/key-locker/proto/keylocker"
)

const ChainName = "Ipfs"

type KeyAdaptor struct {
	fallback.KeyAdaptor
	clients *multiclient.MultiClient
}

func NewChainAdaptor(conf *config.Config) (blockchain.KeyAdaptor, error) {
	return &KeyAdaptor{
		clients: nil,
	}, nil
}

func NewLocalKeyAdaptor(network config.NetWorkType) blockchain.KeyAdaptor {
	return newKeyAdaptor()
}

func newKeyAdaptor() blockchain.KeyAdaptor {
	return &KeyAdaptor{
		clients: multiclient.New([]multiclient.Client{}),
	}
}

func (a *KeyAdaptor) getClient() {

}

func (a *KeyAdaptor) GetSupportChain(req *keylocker.SupportChainReq) (*keylocker.SupportChainRep, error) {
	return &keylocker.SupportChainRep{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get support chain success",
	}, nil
}

func (a *KeyAdaptor) GetSocialKey(req *keylocker.GetSocialKeyReq) (*keylocker.GetSocialKeyRep, error) {
	// todo: 调用 IPFS 接口获取 Social Key
	// 1. req.uuid 取到链上的存储 req.key 的文件，
	// 2. 解密返回就行
	return &keylocker.GetSocialKeyRep{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get ipfs social key success",
	}, nil
}

func (a *KeyAdaptor) SetSocialKey(req *keylocker.SetSocialKeyReq) (*keylocker.SetSocialKeyRep, error) {
	// todo: 调用 IPFS 接口上传 Social Key
	// 1. 如果对应 uuid(req.uuid) 没有 rsa 密钥对，生成 rsa 密钥对，生成 RSA 密钥对处理，私钥，用用户密码(req.password)进行 AES 加密存储， 公钥明文存储, 如果有直接使用
	// 2. 用 rsa 私钥对 key(req.key 是用户上传的一个私钥) 加密，加密 key 调用 ipfs 上传
	// 3. 返回加密的 RSA 的私钥和明文的 RSA 公钥匙, 加密方式，IPFS 对应的 CID
	return &keylocker.SetSocialKeyRep{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "set ipfs social key success",
	}, nil
}
