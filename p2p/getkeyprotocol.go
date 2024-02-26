package p2p

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/ethereum/go-ethereum/log"
	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/savour-labs/key-locker/proto/keylocker"
	protocols_p2p "github.com/savour-labs/key-locker/proto/p2p"
)

const (
	GetKeyRequest  = "/getkeyreq/0.0.1"
	GetKeyResponse = "/getkeyresp/0.0.1"
	v1             = "1.0.0"
)

type GetKeyProtocol struct {
	node     *P2PHost // local host
	mu       sync.Mutex
	Requests map[string]chan *keylocker.GetSocialKeyRep //[]byte //*keylocker.GetSocialKeyReq // used to access request data from response handlers
	done     chan bool                                  // only for demo purposes to hold main from terminating
}

func NewGetKeyProtocol(node *P2PHost) *GetKeyProtocol {
	//e := GetKeyProtocol{node: node, requests: make(map[string]*keylocker.GetSocialKeyReq)}
	e := GetKeyProtocol{node: node, Requests: make(map[string]chan *keylocker.GetSocialKeyRep)}
	node.SetStreamHandler(GetKeyRequest, e.onGetKeyRequest)
	node.SetStreamHandler(GetKeyResponse, e.onGetKeyResponse)

	// design note: to implement fire-and-forget style messages you may just skip specifying a response callback.
	// a fire-and-forget message will just include a request and not specify a response object

	return &e
}

// peer handle remote stream with getKeyReq protocol
func (p *GetKeyProtocol) onGetKeyRequest(s network.Stream) {
	// get wrapped request data
	data := &protocols_p2p.ReylayMsg{}
	buf, err := io.ReadAll(s)
	if err != nil {
		s.Reset()
		log.Error("read from getKeyReq stream", err)
		return
	}
	s.Close()

	// unmarshal it
	err = proto.Unmarshal(buf, data)
	if err != nil {
		log.Error("unmarshal failed", err)
		return
	}

	log.Debug("%s: Received get-key request from %s. MessageID: %s\n", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.GetMessageData().Id)

	if !p.node.authenticateMessage(data, data.GetMessageData()) {
		log.Error("Failed to authenticate message")
		return
	}

	getKeyReq := &keylocker.GetSocialKeyReq{}
	if err := proto.Unmarshal(data.GetMessageData().Data, getKeyReq); err != nil {
		log.Error("unmarshal keylocker.GetSocialKeyReq failed", err)
		return
	}

	// handle ge-key req from other p2p node
	resp, err := p.node.Dispatcher.Registry[getKeyReq.Chain].GetSocialKey(context.Background(), getKeyReq)
	if err != nil {
		log.Error("handle relay get-key req failed", err)
		return
	}
	respBytes, err := proto.Marshal(resp)
	if err != nil {

	}

	relayResp := &protocols_p2p.ReylayMsg{
		MessageData: p.node.NewMessageData(data.MessageData.Id, false),
		Version:     v1,
	}
	relayResp.MessageData.Data = respBytes
	relayResp.MessageData.MsgType = "get-key response type"
	relayResp.MessageData.Id = data.GetMessageData().GetId()

	// sign the data
	signature, err := p.node.signProtoMessage(relayResp)
	if err != nil {
		log.Error("failed to sign response", err)
		return
	}

	// add the signature to the message
	relayResp.MessageData.Sign = signature

	// send the response
	ok := p.node.sendProtoMessage(s.Conn().RemotePeer(), GetKeyResponse, resp)

	if ok {
		log.Info("get-key response sent to %s from %s.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
	}
	p.done <- true // todo: defer
}

// peer handle remote stream with getKeyResp protocol
func (p *GetKeyProtocol) onGetKeyResponse(s network.Stream) {

	// get wrapped request data and verify
	relayMsg, err := VerifyStreamMsg(s, p)
	if err != nil {
		log.Error("verify stream msg failed", err)
		return
	}

	respCh, ok := p.Requests[relayMsg.GetMessageData().GetId()]
	if !ok {
		log.Error(fmt.Sprintf("found request: %d respCh failed", relayMsg.GetMessageData().GetId()))
		return
	}

	resp := &keylocker.GetSocialKeyRep{}
	if err = proto.Unmarshal(relayMsg.MessageData.Data, resp); err != nil {
		log.Error("unmarshal relay msg response failed", err)
		respCh <- nil
		return
	}
	respCh <- resp
}

func VerifyStreamMsg(s network.Stream, p *GetKeyProtocol) (*protocols_p2p.ReylayMsg, error) {
	// get wrapped request data
	data := &protocols_p2p.ReylayMsg{}
	buf, err := io.ReadAll(s)
	if err != nil {
		s.Reset()
		log.Error("read from getKeyReq stream", err)
		return nil, err
	}
	s.Close()

	// unmarshal it
	err = proto.Unmarshal(buf, data)
	if err != nil {
		log.Error("unmarshal failed", err)
		return nil, err
	}

	log.Debug("%s: Received get-key request from %s. MessageID: %s\n", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.GetMessageData().Id)

	if !p.node.authenticateMessage(data, data.GetMessageData()) {
		log.Error("Failed to authenticate message")
		return nil, err
	}
	return data, nil
}
