package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	_ "github.com/venomuz/service-apiswag-post-user/API-GATEWAY/api/model"
	pb "github.com/venomuz/service-apiswag-post-user/API-GATEWAY/genproto"
	l "github.com/venomuz/service-apiswag-post-user/API-GATEWAY/pkg/logger"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"time"
)

// CreateUser creates users
// @Summary      Create an account
// @Description  This api is for creating user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param 		 user  body model.Useri true "user body"
// @Success      200  {string}  Ok
// @Router       /v1/users [post]
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
	)
	body := pb.Useri{}
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().Create(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create user", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetUser gets user by id
// route /v1/users/{id} [get]
func (h *handlerV1) GetUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().GetByID(
		ctx, &pb.GetIdFromUser{Id: guid})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}
