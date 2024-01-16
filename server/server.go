package server

import (
	"GoGinBoilerPlate/server/api"
	"GoGinBoilerPlate/server/props"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Server struct {
	api         api.GoGinHandler
	contextRoot string
	doOnce      sync.Once
	host        string
	port        int
}

func NewServer(properties *props.Properties) *Server {
	s := new(Server)
	s.host = properties.Server.Host
	s.port = properties.Server.Port
	s.contextRoot = properties.Server.ContextRoot
	s.api = api.NewGoGin()
	return s
}

// ConfigureAPI configures the API with all the endpoints with respective handlers
func (s *Server) ConfigureAPI(goGinDomain service.GoGinDomain) {
	s.doOnce.Do(func() {
		s.api = configureAPI(s.api, s.contextRoot, goGinDomain)
	})
}

// Serve the api
func (s *Server) Serve() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
		Handler: s.api,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Fatalf("listen: %s\n", err)
		}
	}()

	s.api.Logger(context.Background(), "Started server on : %s", fmt.Sprintf("http://%s:%d", s.host, s.port))
	//s.api.Logger("Started server on : %s", fmt.Sprintf("https://%s:%d", s.host, s.port))

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.api.Logger(context.Background(), "Server is shutting down ... %s", "quiz-domain-handler")

	s.Shutdown()
	//  ctx, cancel := context.WithTimeout(context.Background(), s.cleanupTimeout)
	// defer cancel()
	// if err := srv.Shutdown(ctx); err != nil {
	// 	s.Logf("Server Shutdown: %v", err)
	// 	return err
	// }

	s.api.Logger(context.Background(), "Until next time ... %s")
	return nil
}

// Fatalf logs message either via defined user logger or via system one if no user logger is defined.
// Exits with non-zero status after printing
func (s *Server) Fatalf(f string, args ...interface{}) {
	if s.api != nil {
		s.api.Logger(context.Background(), f, args...)
		os.Exit(1)
	} else {
		log.Fatalf(f, args...)
	}
}

// Logf logs message either via defined user logger or via system one if no user logger is defined.
func (s *Server) Logf(f string, args ...interface{}) {
	if s.api != nil {
		s.api.Logger(context.Background(), f, args...)
		//s.api.Logger(f, args...)
	} else {
		log.Printf(f, args...)
	}
}

// SetAPI configures the server with the specified API. Needs to be called before Serve
func (s *Server) SetAPI(restApi api.GoGinHandler) {
	if restApi == nil {
		s.api = nil
		return
	}

	s.api = restApi
}

// Shutdown server and clean up resources
func (s *Server) Shutdown() {
	s.api.ServerShutdown()
}
