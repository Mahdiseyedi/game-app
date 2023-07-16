package presence

import (
	"context"
	"game-app/contract/golang/presence"
	"game-app/param"
	"game-app/pkg/protobufmapper"
	"game-app/pkg/slice"
	"google.golang.org/grpc"
)

type Client struct {
	client presence.PresenceServiceClient
}

func New(conn *grpc.ClientConn) Client {
	return Client{client: presence.NewPresenceServiceClient(conn)}
}

func (c Client) GetPresence(ctx context.Context,
	request param.GetPresenceRequest) (param.GetPresenceResponse, error) {

	resp, err := c.client.GetPresence(ctx, &presence.GetPresenceRequest{
		UserIds: slice.MapFromUintToUint64(request.UserIDs),
	})

	if err != nil {
		return param.GetPresenceResponse{}, err
	}

	return protobufmapper.MapGetPresenceResponseFromProtobuf(resp), nil
}
