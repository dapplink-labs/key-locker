package moonbeam

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/savour-labs/key-locker/blockchain/moonbeam/bindings"
	"github.com/savour-labs/key-locker/config"
	"math/big"
	"strings"
	"time"
)

const (
	UuidSize     = 32
	TaskInterval = 12
)

type KeyLockerClient struct {
	context               context.Context
	walletAddress         common.Address
	ethClient             *ethclient.Client
	klContract            *bindings.KeyLocker
	ChainID               *big.Int
	PrivKey               *ecdsa.PrivateKey
	confirmReceiptTimeout time.Duration
	confirmations         int64
}

func NewKeyLockerClient(conf *config.Config) (*KeyLockerClient, error) {
	chainConfig := params.GoerliChainConfig
	if conf.NetWork == "mainnet" {
		chainConfig = params.MainnetChainConfig
	} else if conf.NetWork == "regtest" {
		chainConfig = params.AllCliqueProtocolChanges
	}
	log.Info("eth client setup", "chain_id", chainConfig.ChainID.Int64(), "network", conf.NetWork)
	client, err := ethclient.Dial(conf.Fullnode.Eth.RPCURL)
	klContract, err := bindings.NewKeyLocker(
		common.HexToAddress(conf.Fullnode.Eth.KeyLockerAddr),
		client,
	)
	if err != nil {
		return nil, err
	}
	hex := strings.TrimPrefix(conf.Fullnode.Eth.WalletPriv, "0x")
	priv, err := crypto.HexToECDSA(hex)
	if err != nil {
		return nil, err
	}
	return &KeyLockerClient{
		context:               context.Background(),
		walletAddress:         common.HexToAddress(conf.Fullnode.Eth.WalletAddr),
		ethClient:             client,
		klContract:            klContract,
		ChainID:               chainConfig.ChainID,
		PrivKey:               priv,
		confirmReceiptTimeout: time.Duration(conf.Fullnode.Eth.TimeOut),
		confirmations:         conf.Fullnode.Eth.Confirmations,
	}, nil
}

func (kl KeyLockerClient) AppendSocialKey(uuid [UuidSize]byte, keys [][]byte) error {
	nonce64, err := kl.ethClient.NonceAt(
		kl.context, kl.walletAddress, nil,
	)
	if err != nil {
		log.Error("can not to get current nonce", "err", err)
		return err
	}
	nonce := new(big.Int).SetUint64(nonce64)
	gasPrice, err := kl.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Error("cannot fetch gas price")
		return err
	}
	opts, err := bind.NewKeyedTransactorWithChainID(
		kl.PrivKey, kl.ChainID,
	)
	opts.Context = kl.context
	opts.Nonce = nonce
	opts.NoSend = true
	opts.GasPrice = gasPrice
	tx, err := kl.klContract.SetSocialKey(opts, uuid, keys)
	if err != nil {
		log.Error("can not to set social key")
		return err
	}
	if err := kl.ethClient.SendTransaction(kl.context, tx); err != nil {
		log.Error("can not to send transaction to l1 chain")
		return err
	}
	confirmTxReceipt := func(txHash common.Hash) *types.Receipt {
		ctx, cancel := context.WithTimeout(context.Background(), kl.confirmReceiptTimeout)
		queryTicker := time.NewTicker(TaskInterval)
		defer func() {
			cancel()
			queryTicker.Stop()
		}()
		for {
			receipt, err := kl.ethClient.TransactionReceipt(context.Background(), txHash)
			switch {
			case receipt != nil:
				txHeight := receipt.BlockNumber.Uint64()
				tipHeight, err := kl.ethClient.BlockNumber(context.Background())
				if err != nil {
					log.Error("can not to fetch block number", "err", err)
					break
				}
				log.Info("Transaction mined, checking confirmations",
					"txHash", txHash, "txHeight", txHeight,
					"tipHeight", tipHeight,
					"numConfirmations", kl.confirmations)
				if txHeight+uint64(kl.confirmations) < tipHeight {
					reverted := receipt.Status == 0
					log.Info("Transaction confirmed",
						"txHash", txHash,
						"reverted", reverted)
					return receipt
				}
			case err != nil:
				log.Error("failed to query receipt for transaction", "txHash", txHash.String())
			default:
			}
			select {
			case <-ctx.Done():
				return nil
			case <-queryTicker.C:
			}
		}
	}
	go confirmTxReceipt(tx.Hash())
	if err != nil {
		return err
	}
	return nil

}

func (kl KeyLockerClient) QuerySocialKey(uuid [UuidSize]byte) ([][]byte, error) {
	keys, err := kl.klContract.GetSocialKey(&bind.CallOpts{
		Pending: false,
		Context: kl.context,
	}, uuid)
	if err != nil {
		log.Error("can not to get social key")
		return nil, err
	}
	return keys, nil
}
