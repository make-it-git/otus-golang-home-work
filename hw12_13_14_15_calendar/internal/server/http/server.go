package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	logger2 "github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/logic"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/config"
)

type Server struct {
	config *config.HTTPConf
	logger logger2.ILogger
	server *http.Server
	app    *app.App
}

type event struct {
	ID               string     `json:"id"`
	Title            string     `json:"title"`
	StartTime        time.Time  `json:"startTime"`
	EndTime          time.Time  `json:"endTime"`
	Description      *string    `json:"description,omitempty"`
	OwnerID          int32      `json:"ownerId"`
	NotificationTime *time.Time `json:"notificationTime,omitempty"`
}

func eventToStorage(e *event) storage.Event {
	return storage.Event{
		ID:               e.ID,
		Title:            e.Title,
		StartTime:        e.StartTime,
		Duration:         e.EndTime.Sub(e.StartTime),
		Description:      e.Description,
		OwnerID:          e.OwnerID,
		NotificationTime: e.NotificationTime,
	}
}

func storageToEvent(ev *storage.Event) *event {
	return &event{
		ID:               ev.ID,
		Title:            ev.Title,
		StartTime:        ev.StartTime,
		EndTime:          ev.StartTime.Add(ev.Duration),
		Description:      ev.Description,
		OwnerID:          ev.OwnerID,
		NotificationTime: ev.NotificationTime,
	}
}

func storageToEventList(events []storage.Event) []*event {
	r := make([]*event, 0, len(events))
	for _, ev := range events {
		r = append(r, storageToEvent(&ev))
	}
	return r
}

func writeData(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

func (s *Server) create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	var e event
	err = json.Unmarshal(body, &e)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	ev := eventToStorage(&e)
	err = s.app.CreateEvent(r.Context(), ev)
	if err != nil {
		if errors.As(err, &logic.ErrBussinessLogic{}) {
			writeError(w, http.StatusBadRequest, err)
		} else {
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	writeData(w, http.StatusCreated, storageToEvent(&ev))
}

func (s *Server) update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	var e event
	err = json.Unmarshal(body, &e)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	e.ID = p.ByName("id")
	ev := eventToStorage(&e)
	err = s.app.UpdateEvent(r.Context(), ev)
	if err != nil {
		if errors.As(err, &logic.ErrBussinessLogic{}) {
			writeError(w, http.StatusBadRequest, err)
		} else {
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	writeData(w, http.StatusOK, storageToEvent(&ev))
}

func (s *Server) delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	err := s.app.DeleteEvent(r.Context(), id)
	if err != nil {
		if errors.Is(err, logic.ErrNotFoundID) {
			writeError(w, http.StatusNotFound, err)
		} else if errors.As(err, &logic.ErrBussinessLogic{}) {
			writeError(w, http.StatusBadRequest, err)
		} else {
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	writeData(w, http.StatusNoContent, nil)
}

func (s *Server) listDay(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	date, err := time.Parse("2006-01-02", p.ByName("date"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	events, err := s.app.ListDay(r.Context(), date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeData(w, http.StatusOK, storageToEventList(events))
}

func (s *Server) listWeek(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	date, err := time.Parse("2006-01-02", p.ByName("date"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	events, err := s.app.ListWeek(r.Context(), date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeData(w, http.StatusOK, storageToEventList(events))
}

func (s *Server) listMonth(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	date, err := time.Parse("2006-01-02", p.ByName("date"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	events, err := s.app.ListMonth(r.Context(), date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeData(w, http.StatusOK, storageToEventList(events))
}

func NewServer(logger logger2.ILogger, config *config.HTTPConf, app *app.App) *Server {
	return &Server{
		logger: logger,
		config: config,
		app:    app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	addr := net.JoinHostPort(s.config.Host, s.config.Port)
	s.logger.Info("Start http listen at " + addr)
	router := s.setupRouter()
	configuredRouter := loggingMiddleware(s.logger)(router)
	s.server = &http.Server{Addr: addr, Handler: configuredRouter}
	err := s.server.ListenAndServe()
	if err != nil {
		s.logger.Error("Failed listen at " + addr)
		return err
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}

func (s *Server) setupRouter() *httprouter.Router {
	router := httprouter.New()
	router.POST("/events/", s.create)
	router.PUT("/events/:id", s.update)
	router.DELETE("/events/:id", s.delete)
	router.GET("/events/day/:date", s.listDay)
	router.GET("/events/week/:date", s.listWeek)
	router.GET("/events/month/:date", s.listMonth)
	return router
}
