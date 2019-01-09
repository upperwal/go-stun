package stun

import (
	"net"
	"sync"
	"time"

	"github.com/multiformats/go-multiaddr-net"

	proto "github.com/gogo/protobuf/proto"
	logging "github.com/ipfs/go-log"
	ic "github.com/libp2p/go-libp2p-crypto"
	quic "github.com/lucas-clemente/quic-go"
	ma "github.com/multiformats/go-multiaddr"
	protocol "github.com/upperwal/go-stun/protocol"
)

var log = logging.Logger("stun")

var quicConfig = &quic.Config{
	Versions:                              []quic.VersionNumber{quic.VersionMilestone0_10_0},
	MaxIncomingStreams:                    1000,
	MaxIncomingUniStreams:                 -1,              // disable unidirectional streams
	MaxReceiveStreamFlowControlWindow:     3 * (1 << 20),   // 3 MB
	MaxReceiveConnectionFlowControlWindow: 4.5 * (1 << 20), // 4.5 MB
	AcceptCookie: func(clientAddr net.Addr, cookie *quic.Cookie) bool {
		// TODO(#6): require source address validation when under load
		return true
	},
	KeepAlive:   true,
	IdleTimeout: 30 * time.Hour,
}

type ProtocolPacket struct {
	protocol *protocol.Stun
	raddr    net.Addr
}

type ServerOptions struct {
}

type Server struct {
	conn      net.PacketConn
	l         quic.Listener
	streamMap map[string]*quic.Stream
	mapMutex  *sync.Mutex
}

func NewServer(key ic.PrivKey) (*Server, error) {
	tlsConfig, err := generateConfig(key)
	if err != nil {
		return nil, err
	}

	c, err := net.ListenPacket("udp4", "0.0.0.0:3000")
	if err != nil {
		return nil, err
	}
	log.Info("Listening on: ", c.LocalAddr())

	l, err := quic.Listen(c, tlsConfig, quicConfig)
	if err != nil {
		return nil, err
	}

	s := &Server{
		conn:      c,
		l:         l,
		streamMap: make(map[string]*quic.Stream),
		mapMutex:  &sync.Mutex{},
	}

	go func() {
		for {
			sess, err := s.l.Accept()
			if err != nil {
				log.Error(err)
				return
			}
			log.Info("New Session from: ", sess.RemoteAddr())
			go s.run(sess)
		}
	}()

	return s, nil
}

func (s *Server) run(sess quic.Session) {
	defer sess.Close()
	unmarshalData := &protocol.Stun{}
	buf := make([]byte, 1000)

	log.Info("Waiting for a stream...")
	stream, err := sess.AcceptStream()
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("New stream from: ", sess.RemoteAddr())

	raddr := sess.RemoteAddr()
	log.Info("Connection raddr: ", raddr.String())
	s.mapMutex.Lock()
	s.streamMap[raddr.String()] = &stream
	s.mapMutex.Unlock()

	for {
		i, err := stream.Read(buf)
		if err != nil {
			log.Error(err)
			return
		}

		proto.Unmarshal(buf[:i], unmarshalData)

		pp := ProtocolPacket{
			protocol: unmarshalData,
			raddr:    sess.RemoteAddr(),
		}

		switch unmarshalData.Type {
		case protocol.Stun_HOLE_PUNCH_REQUEST:
			go s.handleHolePunchRequest(pp)
		case protocol.Stun_HOLE_PUNCH_REQUEST_ACCEPT:
			go s.handleHolePunchRequestAccept(pp)
		case protocol.Stun_KEEP_ALIVE:
			// Ignore keep alive messages
		case protocol.Stun_CONNECT:
			log.Info("Got connect packet from: ", raddr)
		}
	}
}

func (s *Server) handleHolePunchRequest(pp ProtocolPacket) {
	peerma, err := ma.NewMultiaddrBytes(pp.protocol.HolePunchRequestMessage.ConnectToPeerID)
	if err != nil {
		log.Error(err)
		return
	}
	peerAddr, err := manet.ToNetAddr(peerma)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(peerAddr.String())
	foreignStream, ok := s.streamMap[peerAddr.String()]
	if !ok {
		log.Error("No stream to this peer")
		return
	}

	raddr, err := manet.FromNetAddr(pp.raddr)
	if err != nil {
		log.Error(err)
	}
	log.Info("New hole punch: foreign -> ", peerAddr)

	packet := &protocol.Stun{
		Type: protocol.Stun_HOLE_PUNCH_REQUEST,
		HolePunchRequestMessage: &protocol.Stun_HolePunchRequestMessage{
			ConnectToPeerID: raddr.Bytes(),
		},
	}

	raw, err := proto.Marshal(packet)
	if err != nil {
		log.Error(err)
	}

	log.Info("Sending signal to foreign peer")
	_, err = (*foreignStream).Write(raw)
	if err != nil {
		log.Error(err)
	}
}

func (s *Server) handleHolePunchRequestAccept(pp ProtocolPacket) {
	peerma, err := ma.NewMultiaddrBytes(pp.protocol.HolePunchRequestMessage.ConnectToPeerID)
	if err != nil {
		log.Error(err)
		return
	}
	peerAddr, err := manet.ToNetAddr(peerma)
	if err != nil {
		log.Error(err)
		return
	}
	initiatorStream, ok := s.streamMap[peerAddr.String()]
	if !ok {
		log.Error("No stream to this peer")
		return
	}

	raddr, err := manet.FromNetAddr(pp.raddr)
	if err != nil {
		log.Error(err)
	}

	log.Info("Hole punch acceptance: initiator -> ", peerAddr)

	packet := &protocol.Stun{
		Type: protocol.Stun_HOLE_PUNCH_REQUEST_ACCEPT,
		HolePunchRequestMessage: &protocol.Stun_HolePunchRequestMessage{
			ConnectToPeerID: raddr.Bytes(),
		},
	}

	raw, err := proto.Marshal(packet)
	if err != nil {
		log.Error(err)
	}

	log.Info("Sending signal to initiator peer")
	_, err = (*initiatorStream).Write(raw)
	if err != nil {
		log.Error(err)
	}

}
