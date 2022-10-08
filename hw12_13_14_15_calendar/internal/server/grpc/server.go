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
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"net"
	"time"
)

type Server struct {
	eventpb.UnimplementedEventServiceServer
	server *grpc.Server
	logger *logger.Logger
	config *config.GRPCConf
	app    *app.App
}

func NewServer(logger *logger.Logger, config *config.GRPCConf, app *app.App) *Server {
	grpcServer := grpc.NewServer(
		withServerUnaryInterceptor(logger),
	)
	s := &Server{
		server: grpcServer,
		logger: logger,
		config: config,
		app:    app,
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
	if s.server != nil {
		s.server.GracefulStop()
	}

	return nil
}

func pbToStorage(in *eventpb.Event) *storage.Event {
	var desc *string
	if in.Description.GetValue() != "" {
		d := in.Description.GetValue()
		desc = &d
	}
	var nt *time.Time
	if in.NotificationTime.IsValid() {
		t := in.NotificationTime.AsTime()
		nt = &t
	}
	return &storage.Event{
		ID:               in.Id,
		Title:            in.Title,
		StartTime:        in.StartTime.AsTime(),
		Duration:         in.EndTime.AsTime().Sub(in.StartTime.AsTime()),
		Description:      desc,
		OwnerID:          in.OwnerId,
		NotificationTime: nt,
	}
}

func storageToPb(ev *storage.Event) *eventpb.Event {
	var ntpb *timestamppb.Timestamp
	var descpb *wrapperspb.StringValue
	if ev.NotificationTime != nil {
		ntpb = timestamppb.New(*ev.NotificationTime)
	}
	if ev.Description != nil {
		descpb = wrapperspb.String(*ev.Description)
	}
	return &eventpb.Event{
		Id:               ev.ID,
		Title:            ev.Title,
		StartTime:        timestamppb.New(ev.StartTime),
		EndTime:          timestamppb.New(ev.StartTime.Add(ev.Duration)),
		Description:      descpb,
		OwnerId:          ev.OwnerID,
		NotificationTime: ntpb,
	}
}

func storageToPbList(events []storage.Event) *eventpb.EventList {
	e := new(eventpb.EventList)
	e.Events = make([]*eventpb.Event, 0, len(events))
	for _, ev := range events { //nolint:typecheck
		e.Events = append(e.Events, storageToPb(&ev))
	}
	return e
}

func (s *Server) Create(ctx context.Context, in *eventpb.Event) (*eventpb.Event, error) {
	ev := pbToStorage(in)
	err := s.app.CreateEvent(ctx, *ev)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return storageToPb(ev), nil
}
func (s *Server) Update(ctx context.Context, in *eventpb.Event) (*eventpb.Event, error) {
	ev := pbToStorage(in)
	err := s.app.UpdateEvent(ctx, *ev)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return storageToPb(ev), nil
}
func (s *Server) Delete(ctx context.Context, in *eventpb.EventId) (*emptypb.Empty, error) {
	err := s.app.DeleteEvent(ctx, in.Id)
	return nil, err
}
func (s *Server) ListDay(ctx context.Context, in *eventpb.EventDate) (*eventpb.EventList, error) {
	events, err := s.app.ListDay(ctx, in.Date.AsTime())
	if err != nil {
		return nil, err
	}
	return storageToPbList(events), nil
}
func (s *Server) ListWeek(ctx context.Context, in *eventpb.EventDate) (*eventpb.EventList, error) {
	events, err := s.app.ListWeek(ctx, in.Date.AsTime())
	if err != nil {
		return nil, err
	}
	return storageToPbList(events), nil
}
func (s *Server) ListMonth(ctx context.Context, in *eventpb.EventDate) (*eventpb.EventList, error) {
	events, err := s.app.ListMonth(ctx, in.Date.AsTime())
	if err != nil {
		return nil, err
	}
	return storageToPbList(events), nil
}
