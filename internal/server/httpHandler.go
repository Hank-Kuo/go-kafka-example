package server

import (
	"net/http"
	"time"

	userDelivery "github.com/Hank-Kuo/go-kafka-example/internal/api/delivery/user"
	userRepository "github.com/Hank-Kuo/go-kafka-example/internal/api/repository/user"
	userService "github.com/Hank-Kuo/go-kafka-example/internal/api/service/user"
	http_middleware "github.com/Hank-Kuo/go-kafka-example/internal/middleware/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) registerHttpHanders_() {
	api := s.engine.Group("/api")

	userRepo := userRepository.NewRepo(s.db, s.kakfaWriter)
	userSrv := userService.NewService(s.cfg, userRepo, s.kakfaWriter, s.logger)
	userDelivery.NewHttpHandler(api, userSrv, s.logger)
}

func (s *Server) registerHttpHanders(engine *gin.Engine) {
	api := engine.Group("/api")

	userRepo := userRepository.NewRepo(s.db, s.kakfaWriter)
	userSrv := userService.NewService(s.cfg, userRepo, s.kakfaWriter, s.logger)
	userDelivery.NewHttpHandler(api, userSrv, s.logger)
}

func (s *Server) newHttpServer() *http.Server {
	if s.cfg.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()
	http_middleware.NewGlobalMiddlewares(engine)

	s.registerHttpHanders(engine)

	httpServer := &http.Server{
		Addr:           ":" + s.cfg.Server.Port,
		Handler:        engine,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Error http ListenAndServe: %s", err)
		}
	}()

	return httpServer
}
