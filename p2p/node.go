package p2p

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/log"
	ggio "github.com/gogo/protobuf/io"
	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/savour-labs/key-locker/proto/keylocker"
	protocols_p2p "github.com/savour-labs/key-locker/proto/p2p"
	"golang.org/x/exp/rand"
)

const (
	clientVersion = "go-p2p-node/0.0.1"
	directReqType = "directFromCliType"
	relayReqType  = "relayFromP2PType"
)

// helper method - generate message data shared between all node's p2p protocols
// messageId: unique for requests, copied from request for responses
func (n *P2PHost) NewMessageData(messageId string, gossip bool) *protocols_p2p.MessageData {
	// Add protobuf bin data for message author public key
	// this is useful for authenticating  messages forwarded by a node authored by another node

	//nodePubKey, err := crypto.MarshalPublicKey(n.Peerstore().PubKey(n.ID()))
	//
	//if err != nil {
	//	panic("Failed to get public key for sender from local peer store.")
	//}

	return &protocols_p2p.MessageData{ClientVersion: clientVersion,
		NodeId: n.ID().String(),
		//NodePubKey: nodePubKey,
		Timestamp: time.Now().Unix(),
		Id:        messageId,
		Gossip:    gossip}
}

// FaultTolerance is a general fault tolerance function
// turn the req handled failed locally to p2p peer according to p2pPeerStore
// know originReq type from protocol
// func (h *P2PHost) FaultTolerance(protocol protocol.ID, originReq *keylocker.GetSocialKeyReq) ([]byte, error) {
func (h *P2PHost) FaultTolerance(protocol protocol.ID, originReqBytes []byte) (string, error) {
	ctx := context.Background()
	peerNum := h.Peerstore().Peers().Len()
	if peerNum == 0 {
		return "", errors.New("no peers in peerStore")
	}
	if originReqBytes == nil {
		return "", errors.New("param req with nil pointer")
	}
	//originReqBytes, err := proto.Marshal(originReq)
	//if err != nil {
	//	log.Error("proto marshal origin request failed", err)
	//	return "", err
	//}

	req := &protocols_p2p.ReylayMsg{
		MessageData: h.NewMessageData(uuid.New().String(), false), // UUID是一个128位的标识符，通常用于唯一标识实体，以确保在分布式系统中生成的标识符不会发生冲突。
		Version:     v1,
	}
	req.MessageData.MsgType = relayReqType
	req.MessageData.Data = originReqBytes
	// sign req.hash signature
	signature, err := h.signProtoMessage(req.MessageData)
	if err != nil {
		log.Error("failed to sign pb data")
		return "", err
	}
	// add the signature to the req message
	req.MessageData.Sign = signature

	// store ref request so response handler has access to it
	h.GetKeyProtocol.mu.Lock()
	h.GetKeyProtocol.Requests[req.MessageData.Id] = make(chan *keylocker.GetSocialKeyRep, 1) //originReqBytes
	h.GetKeyProtocol.mu.Unlock()

	// define retry time
	var retryTimes int
	if peerNum <= h.FaultToleranceTimes {
		retryTimes = peerNum
	} else {
		retryTimes = h.FaultToleranceTimes
	}
	// relay
	for ; retryTimes > 0; retryTimes-- {
		peers := h.Peerstore().PeersWithAddrs()
		rand.Seed(uint64(time.Now().Nanosecond()))
		p2pPeerInfo := h.Peerstore().PeerInfo(peers[rand.Intn(peerNum)])
		s, err := h.NewStream(ctx, p2pPeerInfo.ID, protocol)
		if err != nil {
			continue
		}
		defer s.Close()

		writer := ggio.NewFullWriter(s)
		err = writer.WriteMsg(req)
		if err == nil {
			return req.MessageData.Id, nil
		}
		s.Reset()
	}
	return "", errors.New("found no p2p host to relay req")
}

// sign an outgoing p2p message payload
func (n *P2PHost) signProtoMessage(message proto.Message) ([]byte, error) {
	data, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	return n.signData(data)
}

// sign binary data using the local node's private key
func (n *P2PHost) signData(data []byte) ([]byte, error) {
	key := n.Peerstore().PrivKey(n.ID())
	res, err := key.Sign(data)
	return res, err
}

// Authenticate incoming p2p message
// message: a protobuf go data object
// data: common p2p message data
func (n *P2PHost) authenticateMessage(message proto.Message, data *protocols_p2p.MessageData) bool {
	// store a temp ref to signature and remove it from message data
	// sign is a string to allow easy reset to zero-value (empty string)
	signTmp := data.Sign
	data.Sign = nil

	// marshall data without the signature to protobufs3 binary format
	bin, err := proto.Marshal(message)
	if err != nil {
		log.Error("failed to marshal pb message", err)
		return false
	}

	// restore sig in message data (for possible future use)
	data.Sign = signTmp

	// restore peer id binary format from base58 encoded node id data
	peerId, err := peer.Decode(data.NodeId)
	if err != nil {
		log.Error("Failed to decode node id from base58", err)
		return false
	}

	// verify the data was authored by the signing peer identified by the public key
	// and signature included in the message
	return n.verifyData(bin, signTmp, peerId, data.NodePubKey)
}

// Verify incoming p2p message data integrity
// data: data to verify
// signature: author signature provided in the message payload
// peerId: author peer id from the message payload
// pubKeyData: author public key from the message payload
func (n *P2PHost) verifyData(data []byte, signature []byte, peerId peer.ID, pubKeyData []byte) bool {
	//key, err := crypto.UnmarshalPublicKey(pubKeyData)
	//if err != nil {
	//	log.Println(err, "Failed to extract key from message key data")
	//	return false
	//}
	//
	//// extract node id from the provided public key
	//idFromKey, err := peer.IDFromPublicKey(key)
	//
	//if err != nil {
	//	log.Println(err, "Failed to extract peer id from public key")
	//	return false
	//}
	//
	//// verify that message author node id matches the provided node public key
	//if idFromKey != peerId {
	//	log.Println(err, "Node id and provided public key mismatch")
	//	return false
	//}
	key := n.Peerstore().PubKey(peerId)

	res, err := key.Verify(data, signature)
	if err != nil {
		log.Error("Error authenticating data", err)
		return false
	}

	return res
}

// helper method - writes a protobuf go data object to a network stream
// data: reference of protobuf go data object to send (not the object itself)
// s: network stream to write the data to
func (n *P2PHost) sendProtoMessage(id peer.ID, p protocol.ID, data proto.Message) bool {
	s, err := n.NewStream(context.Background(), id, p)
	if err != nil {
		log.Error("new stream failed", err)
		return false
	}
	defer s.Close()

	writer := ggio.NewFullWriter(s)
	err = writer.WriteMsg(data)
	if err != nil {
		log.Error("write msg to stream failed", err)
		s.Reset()
		return false
	}
	return true
}
