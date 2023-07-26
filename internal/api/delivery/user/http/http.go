package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	userSrv "go-kafka-example/internal/api/service/user"
	"go-kafka-example/pkg/logger"
	"go-kafka-example/pkg/response"
	"go-kafka-example/pkg/tracer"
)

type httpHandler struct {
	userSrv userSrv.Service
	logger  logger.Logger
}

func NewHandler(e *gin.RouterGroup, userSrv userSrv.Service, logger logger.Logger) {
	handler := &httpHandler{
		userSrv: userSrv,
		logger:  logger,
	}

	e.POST("/register", handler.Create)
	// e.POST("/login", handler.Create)
	e.GET("/users", handler.Get)
}

func (h *httpHandler) Create(c *gin.Context) {

	data := "string"
	response.OK(http.StatusOK, "login successfully", data).ToJSON(c)
}

func (h *httpHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := tracer.NewSpan(ctx, "UserHandler.Login", nil)
	defer span.End()

	data, err := h.userSrv.GetAll(ctx)
	if err != nil {
		tracer.AddSpanError(span, err)
		response.Fail(err, h.logger).ToJSON(c)
		return
	}

	response.OK(http.StatusOK, "login successfully", data).ToJSON(c)
}
