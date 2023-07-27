package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	userSrv "go-kafka-example/internal/api/service/user"
	"go-kafka-example/internal/dto"
	"go-kafka-example/internal/models"
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

	e.POST("/register", handler.Register)
	e.POST("/login", handler.Login)
	e.POST("/user/active", handler.UserActive)
	e.GET("/users", handler.GetUsers)
}

func (h *httpHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	_, span := tracer.NewSpan(ctx, "UserHandler.Register", nil)
	defer span.End()

	var body dto.RegisterReqDto
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(err, h.logger).ToJSON(c)
		return
	}

	err := h.userSrv.Register(ctx, models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
	})

	if err != nil {
		tracer.AddSpanError(span, err)
		response.Fail(err, h.logger).ToJSON(c)
		return
	}

	response.OK(http.StatusOK, "register successfully", nil).ToJSON(c)
}

func (h *httpHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := tracer.NewSpan(ctx, "UserHandler.Login", nil)
	defer span.End()

	var body dto.LoginReqDto
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(err, h.logger).ToJSON(c)
		return
	}

	data, err := h.userSrv.Login(ctx, "test", "password")
	if err != nil {
		tracer.AddSpanError(span, err)
		response.Fail(err, h.logger).ToJSON(c)
		return
	}

	response.OK(http.StatusOK, "login successfully", data).ToJSON(c)
}

func (h *httpHandler) GetUsers(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := tracer.NewSpan(ctx, "UserHandler.GetUsers", nil)
	defer span.End()

	data, err := h.userSrv.GetAll(ctx)
	if err != nil {
		tracer.AddSpanError(span, err)
		response.Fail(err, h.logger).ToJSON(c)
		return
	}

	response.OK(http.StatusOK, "login successfully", data).ToJSON(c)
}

func (h *httpHandler) UserActive(c *gin.Context) {
	ctx := c.Request.Context()
	_, span := tracer.NewSpan(ctx, "UserHandler.UserActive", nil)
	defer span.End()

	response.OK(http.StatusOK, "login successfully", "data").ToJSON(c)
}
