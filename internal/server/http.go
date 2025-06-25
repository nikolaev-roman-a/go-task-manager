package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/nikolaev-roman-a/go-task-manager/internal/models"
	"github.com/nikolaev-roman-a/go-task-manager/internal/services"
	"go.uber.org/zap"
)

type Server struct {
	logger  *zap.Logger
	server  *http.Server
	service *services.TaskService
}

func NewHTTPServer(service *services.TaskService, logger *zap.Logger) *Server {
	mux := http.NewServeMux()
	s := &Server{
		service: service,
		logger:  logger,
		server: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
	}

	mux.HandleFunc("POST /tasks", s.setHandler)
	mux.HandleFunc("GET /tasks/{id}", s.getHandler)
	mux.HandleFunc("GET /tasks", s.allHandler)
	mux.HandleFunc("DELETE /tasks/{id}", s.delHandler)

	return s
}

func (s *Server) Run() {
	go func() {
		s.logger.Info("Starting HTTP server", zap.String("addr", ":8080"))
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("HTTP server error", zap.Error(err))
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) {
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("HTTP shutdown error", zap.Error(err))
	}
}

func (s *Server) setHandler(w http.ResponseWriter, r *http.Request) {

	task := &models.Task{}

	task, err := s.service.CreateAndRun(r.Context(), task)
	if err != nil {
		s.logger.Error("Set error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (s *Server) getHandler(w http.ResponseWriter, r *http.Request) {
	taskIDstr := r.PathValue("id")
	taskIDuuid, err := uuid.Parse(taskIDstr)
	if err != nil {
		s.logger.Error("Set error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	task, err := s.service.Read(r.Context(), taskIDuuid)
	if err != nil {
		s.logger.Error("Set error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (s *Server) delHandler(w http.ResponseWriter, r *http.Request) {
	taskIDstr := r.PathValue("id")
	taskIDuuid, err := uuid.Parse(taskIDstr)
	if err != nil {
		s.logger.Error("Set error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.service.Cancel(r.Context(), taskIDuuid)
	if err != nil {
		s.logger.Error("Set error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) allHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.service.Search(r.Context())
	if err != nil {
		s.logger.Error("Set error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
