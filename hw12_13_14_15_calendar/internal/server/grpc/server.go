package grpc

import (
	"context"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/adapters/grpc/eventpb"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

type Server struct {
	eventpb.UnimplementedEventServiceServer
	server  *grpc.Server
	logger  *logger.Logger
	config  *config.GRPCConf
	storage app.Storage
}

func NewServer(logger *logger.Logger, config *config.GRPCConf, storage app.Storage) *Server {
	grpcServer := grpc.NewServer()
	s := &Server{
		server:  grpcServer,
		logger:  logger,
		config:  config,
		storage: storage,
	}
	eventpb.RegisterEventServiceServer(grpcServer, s)
	reflection.Register(grpcServer)
	return s
}

func (s *Server) Start(ctx context.Context) error {
	addr := net.JoinHostPort(s.config.Host, s.config.Port)
	s.logger.Info("Start grpc listen at " + addr)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	err = s.server.Serve(listen)
	if err != nil {
		s.logger.Error("Failed listen at " + addr)
		return err
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Server) Create(ctx context.Context, in *eventpb.Event) (*eventpb.Event, error) {
	d := in.Description.GetValue()
	n := in.NotificationTime.AsTime()
	ev := storage.Event{
		ID:               in.Id,
		Title:            in.Title,
		StartTime:        in.StartTime.AsTime(),
		Duration:         in.EndTime.AsTime().Sub(in.StartTime.AsTime()),
		Description:      &d,
		OwnerID:          in.OwnerId,
		NotificationTime: &n,
	}
	err := s.storage.Create(ev)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return &eventpb.Event{
		Id:               ev.ID,
		Title:            ev.Title,
		StartTime:        nil,
		EndTime:          nil,
		Description:      nil,
		OwnerId:          0,
		NotificationTime: nil,
	}, nil
}
func (s *Server) Update(ctx context.Context, in *eventpb.Event) (*eventpb.Event, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (s *Server) Delete(ctx context.Context, in *eventpb.EventId) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (s *Server) ListDay(ctx context.Context, in *eventpb.EventDate) (*eventpb.EventList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDay not implemented")
}
func (s *Server) ListWeek(ctx context.Context, in *eventpb.EventDate) (*eventpb.EventList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListWeek not implemented")
}
func (s *Server) ListMonth(ctx context.Context, in *eventpb.EventDate) (*eventpb.EventList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMonth not implemented")
}
