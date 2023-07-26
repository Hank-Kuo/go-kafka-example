package server

import (
	userHttpDelivery "go-kafka-example/internal/api/delivery/user/http"
	userRepository "go-kafka-example/internal/api/repository/user"
	userService "go-kafka-example/internal/api/service/user"
)

func (s *Server) registerHttpHanders() {
	api := s.engine.Group("/api")

	userRepo := userRepository.NewRepo(s.db)
	userSrv := userService.NewService(s.cfg, userRepo, s.logger)
	userHttpDelivery.NewHandler(api, userSrv, s.logger)

}