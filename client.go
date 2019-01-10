package stun

import (
	"crypto/tls"
	"net"
	"sync"
	"time"

	ic "github.com/libp2p/go-libp2p-crypto"
	"github.com/lucas-clemente/quic-go"

	proto "github.com/gogo/protobuf/proto"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
	protocol "github.com/upperwal/go-stun/protocol"
)

type ClientOptions struct {
	serverList []ma.Multiaddr
}

type Client struct {
	conn              net.PacketConn
	stunServerList    []ma.Multiaddr
	stunSession       quic.Session
	stunStream        quic.Stream
	tlsConfig         *tls.Config
	completionMapChan map[string]chan bool
	cmapMutex         *sync.Mutex
}

func NewClient(key ic.PrivKey, pc net.PacketConn) (*Client, error) {
	tlsConfig, err := GenerateConfig(key)
	if err != nil {
		return nil, err
	}
	log.Info("Client starting on: ", pc.LocalAddr())

	return &Client{
		conn:              pc,
		tlsConfig:         tlsConfig,
		completionMapChan: make(map[string]chan bool),
		cmapMutex:         &sync.Mutex{},
	}, nil
}

func (c *Client) ConnectSTUNServer(m []ma.Multiaddr) {
	c.stunServerList = append(c.stunServerList, m...)

	for _, s := range m {
		// TODO: check the error
		raddr, _ := manet.ToNetAddr(s)
		sess, err := quic.Dial(c.conn, raddr, "blockport.p2p", c.tlsConfig, quicConfig)
		if err == nil {
			c.stunSession = sess
			stream, err := sess.OpenStream()
			if err == nil {
				log.Info("New stream to ", sess.RemoteAddr())
				c.stunStream = stream
				c.connectPacket()
				go c.handleMessages()
				return
			}
		}
	}
}

func (c *Client) connectPacket() {
	log.Info("Writing connect packet")
	packetConnect := &protocol.Stun{
		Type: protocol.Stun_CONNECT,
		HolePunchRequestMessage: &protocol.Stun_HolePunchRequestMessage{
			ConnectToPeerID: []byte(""),
		},
	}
	raw, err := proto.Marshal(packetConnect)
	if err != nil {
		log.Error(err)
		return
	}
	_, err = c.stunStream.Write(raw)
	if err != nil {
		log.Error(err)
		return
	}
}

func (c *Client) handleMessages() {
	buf := make([]byte, 1000)
	packet := &protocol.Stun{}

	for {
		i, err := c.stunStream.Read(buf)
		if err != nil {
			log.Error(err)
		}

		err = proto.Unmarshal(buf[:i], packet)
		if err != nil {
			log.Error(err)
		}
		log.Info("Read new message: ", packet)

		switch packet.Type {
		case protocol.Stun_HOLE_PUNCH_REQUEST:
			go c.handleHolePunchRequest(packet)
		case protocol.Stun_HOLE_PUNCH_REQUEST_ACCEPT:
			go c.handleHolePunchRequestAccept(packet)
		}
	}
}

func (c *Client) handleHolePunchRequest(packet *protocol.Stun) {
	log.Info("got a new hole punch request")
	packetAccept := &protocol.Stun{
		Type: protocol.Stun_HOLE_PUNCH_REQUEST_ACCEPT,
		HolePunchRequestMessage: &protocol.Stun_HolePunchRequestMessage{
			ConnectToPeerID: packet.HolePunchRequestMessage.ConnectToPeerID,
		},
	}
	raw, err := proto.Marshal(packetAccept)
	if err != nil {
		return
	}
	_, err = c.stunStream.Write(raw)
	if err != nil {
		return
	}
	log.Info("bombarding now.")
	c.bombardPackets(packet.HolePunchRequestMessage.ConnectToPeerID)
}

func (c *Client) handleHolePunchRequestAccept(packet *protocol.Stun) {
	log.Info("Got hole punching acceptance. bombarding now.")
	c.bombardPackets(packet.HolePunchRequestMessage.ConnectToPeerID)
}

func (c *Client) bombardPackets(peer []byte) {
	maAddr, err := ma.NewMultiaddrBytes(peer)
	if err != nil {
		log.Error(err)
		return
	}
	peerAddr, err := manet.ToNetAddr(maAddr)
	if err != nil {
		log.Error(err)
		return
	}

	//finishChan := make(chan bool)
	//go c.readRawConn(finishChan)

	for i := 0; i < 30; i++ {
		log.Info("Bombarding...", c.conn.LocalAddr())
		_, err := c.conn.WriteTo([]byte("a"), peerAddr)
		if err != nil {
			log.Error(err)
			return
		}
		//time.Sleep(time.Second * 1)
	}

	c.completionMapChan[maAddr.String()] <- true

	//finishChan <- true
}

func (c *Client) readRawConn(f chan bool) {
	buf := make([]byte, 1000)
	log.Info("Reading raw")
	for {
		select {
		case <-f:
			log.Info("Returning finished")
			return
		default:
		}
		i, addr, err := c.conn.ReadFrom(buf)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Info("Reading something:", buf[:i], addr)
	}
}

func (c *Client) PunchHole(raddr ma.Multiaddr) (chan bool, error) {
	log.Info("Punching hole to: ", raddr)
	packet := &protocol.Stun{
		Type: protocol.Stun_HOLE_PUNCH_REQUEST,
		HolePunchRequestMessage: &protocol.Stun_HolePunchRequestMessage{
			ConnectToPeerID: raddr.Bytes(),
		},
	}
	raw, err := proto.Marshal(packet)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	c.cmapMutex.Lock()
	c.completionMapChan[raddr.String()] = make(chan bool)
	c.cmapMutex.Unlock()

	_, err = c.stunStream.Write(raw)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	go func() {
		time.Sleep(time.Second * 1)
		c.completionMapChan[raddr.String()] <- false
	}()

	return c.completionMapChan[raddr.String()], nil
}

/* func (c *Client) sendPunchHoleRequest */
