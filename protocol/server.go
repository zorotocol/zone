package protocol

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/zorotocol/zone/auth"
	"github.com/zorotocol/zone/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"time"
)

var _ pb.ProxyServer = &Server{}

type Server struct {
	pb.UnimplementedProxyServer
	PublicKeys  [][]byte
	DialTimeout time.Duration
	EnableUDP   bool
}

func (s *Server) IDs(context.Context, *empty.Empty) (*pb.IDsResponse, error) {
	return &pb.IDsResponse{Id: s.PublicKeys}, nil
}

func (s *Server) TCP(stream pb.Proxy_TCPServer) error {
	token, dest, err := metadataFromIncomingContext(stream.Context(), true)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if auth.Verify(token, s.PublicKeys) < 0 {
		return status.Error(codes.Unauthenticated, "")
	}
	conn, err := net.DialTimeout("tcp", dest, s.DialTimeout)
	if err != nil {
		return status.Error(codes.Unavailable, "cloud not dial remote")
	}
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	go func() {
		defer cancel()
		buff := make([]byte, 2048)
		chunk := &pb.Chunk{}
		for {
			n, err := conn.Read(buff)
			if err != nil {
				break
			}
			chunk.Data = buff[:n]
			if err = stream.SendMsg(chunk); err != nil {
				break
			}
		}
	}()
	go func() {
		defer cancel()
		chunk := &pb.Chunk{}
		for {
			if err := stream.RecvMsg(chunk); err != nil {
				break
			}
			_, err := conn.Write(chunk.Data)
			if err != nil {
				break
			}
		}
	}()
	<-ctx.Done()
	return nil
}

func (s *Server) UDP(stream pb.Proxy_UDPServer) error {
	if s.EnableUDP == false {
		return status.Error(codes.Unavailable, "")
	}
	token, _, err0 := metadataFromIncomingContext(stream.Context(), false)
	if err0 != nil {
		return status.Error(codes.InvalidArgument, err0.Error())
	}
	if auth.Verify(token, s.PublicKeys) < 0 {
		return status.Error(codes.Unauthenticated, "")
	}
	pconn, err0 := net.ListenPacket("udp", ":0")
	if err0 != nil {
		return status.Error(codes.Internal, "cloud not associate")
	}
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	go func() {
		defer cancel()
		buff := make([]byte, 2048)
		pack := &pb.Packet{}
		for {
			n, from, err := pconn.ReadFrom(buff)
			if err != nil {
				break
			}
			pack.Data = buff[:n]
			pack.Addr = from.String()
			if err = stream.SendMsg(pack); err != nil {
				break
			}
		}
	}()
	go func() {
		var err error
		defer cancel()
		pack := &pb.Packet{}
		for {
			if err = stream.RecvMsg(pack); err != nil {
				break
			}
			udpAddr, ok := resolveUDPAddr(ctx, pack.Addr)
			if !ok {
				break
			}
			_, err = pconn.WriteTo(pack.Data, udpAddr)
			if err != nil {
				break
			}
		}
	}()
	<-ctx.Done()
	return nil
}
