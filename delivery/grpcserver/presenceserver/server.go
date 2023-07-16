package presenceserver

import (
	"context"
	"fmt"
	"game-app/contract/golang/presence"
	"game-app/param"
	"game-app/pkg/protobufmapper"
	"game-app/pkg/slice"
	"game-app/service/presenceservice"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	presence.UnimplementedPresenceServiceServer
	svc presenceservice.Service
}

func New(svc presenceservice.Service) Server {
	return Server{
		UnimplementedPresenceServiceServer: presence.UnimplementedPresenceServiceServer{},
		svc:                                svc,
	}
}

func (s Server) GetPresence(ctx context.Context,
	req *presence.GetPresenceRequest) (*presence.GetPresenceResponse, error) {
	resp, err := s.svc.GetPresence(ctx, param.GetPresenceRequest{
		UserIDs: slice.MapFromUint64ToUint(req.GetUserIds()),
	})

	if err != nil {
		return nil, err
	}

	return protobufmapper.MapGetPresenceResponseToProtobuf(resp), err
}

func (s Server) Start() {

	//step one : like always, create listener on tcp port
	address := fmt.Sprintf(":%d", 8086)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	//step two : pbPresenceserver
	presenceSvcServer := Server{}

	//step three : grpc server
	grpcServer := grpc.NewServer()
	// pbPresenceServer register into grpc server

	presence.RegisterPresenceServiceServer(grpcServer, &presenceSvcServer)
	//server grpcServer by listener

	log.Println("presence.proto grpc server stating on ", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("couldn't serve presence.proto grpc server... ")
	}
}
