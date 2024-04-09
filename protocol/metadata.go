package protocol

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/zorotocol/zone/errorutils"
	"github.com/zorotocol/zone/pb"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func EncodeToken(token *pb.Token) string {
	return base64.RawURLEncoding.EncodeToString(errorutils.Must(proto.Marshal(token)))
}

func SetOutgoingMetadata(md metadata.MD, token, destination string) {
	md["x-destination"] = []string{destination}
	if len(destination) > 0 {
		md["authorization"] = []string{token}
	}
}
func metadataFromIncomingContext(ctx context.Context, needDestination bool) (token *pb.Token, destination string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, "", errors.New("cloud not read metadata")
	}
	tokens, _ := md["authorization"]
	if needDestination {
		destinations, _ := md["x-destination"]
		if len(tokens) != 1 || len(destination) != 1 {
			return nil, "", errors.New("invalid headers")
		}
		tokenBytes, err := base64.RawURLEncoding.DecodeString(tokens[0])
		if err != nil {
			return nil, "", errors.New("invalid token codec")
		}
		token := &pb.Token{}
		return token, destinations[0], proto.Unmarshal(tokenBytes, token)
	} else {
		if len(tokens) != 1 {
			return nil, "", errors.New("invalid headers")
		}
		tokenBytes, err := base64.RawURLEncoding.DecodeString(tokens[0])
		if err != nil {
			return nil, "", errors.New("invalid token codec")
		}
		token := &pb.Token{}
		return token, "", proto.Unmarshal(tokenBytes, token)
	}
}
