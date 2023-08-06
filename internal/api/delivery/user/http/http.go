package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	userSrv "go-kafka-example/internal/api/service/user"
	"go-kafka-example/internal/dto"
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
	e.POST("/send_email", handler.SendEmail)
	e.POST("/login", handler.Login)
	e.POST("/user/active", handler.UserActive)
	e.GET("/users", handler.GetUsers)
}

func (h *httpHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := tracer.NewSpan(ctx, "UserHandler.Register", nil)
	defer span.End()

	var body dto.RegisterReqDto
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(err, h.logger).ToJSON(c)
		return
	}

	if err := h.userSrv.Register(ctx, &body); err != nil {
		tracer.AddSpanError(span, err)
		response.Fail(err, h.logger).ToJSON(c)
		return
	}

	response.OK(http.StatusOK, "register successfully", nil).ToJSON(c)
}

func (h *httpHandler) SendEmail(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := tracer.NewSpan(ctx, "UserHandler.Register", nil)
	defer span.End()

	var body dto.SendEmailReqDto
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(err, h.logger).ToJSON(c)
		return
	}

	user, err := h.userSrv.GetByEmail(ctx, body.Email)
	if err != nil {
		tracer.AddSpanError(span, err)
		response.Fail(err, h.logger).ToJSON(c)
		return
	}

	if user.Status {
		response.OK(http.StatusOK, "you already activate your email", nil).ToJSON(c)
	} else {
		if err := h.userSrv.PublishEmail(ctx, user.Name, user.Email); err != nil {
			tracer.AddSpanError(span, err)
			response.Fail(err, h.logger).ToJSON(c)
			return
		}
		response.OK(http.StatusOK, "send email successfully", nil).ToJSON(c)
	}

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

	data, err := h.userSrv.Login(ctx, body.Email, body.Password)
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

	response.OK(http.StatusOK, "get users successfully", data).ToJSON(c)
}

func (h *httpHandler) UserActive(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := tracer.NewSpan(ctx, "UserHandler.UserActive", nil)
	defer span.End()

	var body dto.UserActiveReqDto
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(err, h.logger).ToJSON(c)
		return
	}
	isActive, err := h.userSrv.Active(ctx, body.Email, body.OTP)
	if err != nil {
		response.Fail(err, h.logger).ToJSON(c)
		return
	}
	response.OK(http.StatusOK, "active user successfully", isActive).ToJSON(c)
}
