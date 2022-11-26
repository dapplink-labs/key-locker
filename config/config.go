package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type RPC struct {
	RPCURL  string `yaml:"rpc_url"`
	RPCUser string `yaml:"rpc_user"`
	RPCPass string `yaml:"rpc_pass"`
}

type Eth struct {
	RPCURL        string `yaml:"rpc_url"`
	RPCUser       string `yaml:"rpc_user"`
	RPCPass       string `yaml:"rpc_pass"`
	KeyLockerAddr string `yaml:"key_locker_addr"`
	WalletAddr    string `yaml:"wallet_addr"`
	WalletPriv    string `yaml:"wallet_priv"`
	TimeOut       int64  `yaml:"time_out"`
	Confirmations int64  `yaml:"confirmations"`
}

// Fullnode define
type Fullnode struct {
	Eth  *Eth  `yaml:"eth"`
	Ipfs *Ipfs `yaml:"ipfs"`
}

type Config struct {
	Database  *Database  `yaml:"database"`
	Fullnode  Fullnode   `yaml:"fullnode"`
	NetWork   string     `yaml:"network"`
	Server    *Server    `yaml:"server"`
	RpcServer *RpcServer `yaml:"rpcserver"`
	Chains    []string   `yaml:"chains"`
	AesKey    string     `yaml:"aes_key"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Server struct {
	Port  int  `yaml:"port"`
	Debug bool `yaml:"debug"`
}

type RpcServer struct {
	Port string `yaml:"port"`
}

type Ipfs struct {
	NetworkNode []string `yaml:"network_node"`
	RepoPath    string   `yaml:"repo_path"`
}

type NetWorkType int

const (
	MainNet NetWorkType = iota
	TestNet
	RegTest
)

func LoadConfigFile(filePath string, cfg *Config) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return yaml.NewDecoder(file).Decode(cfg)
}

const UnsupportedChain = "Unsupport chain"
const UnsupportedOperation = UnsupportedChain
