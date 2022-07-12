package v1

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/venomuz/service-api-swag-gin/ApiGateway/api/model"
	pb "github.com/venomuz/service-api-swag-gin/ApiGateway/genproto"
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
// @Param 		 user  body model.User true "user body"
// @Success      200  {string}  Ok
// @Router       /v1/users [post]
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
	)
	body := pb.User{}
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", logger.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().Create(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create user", logger.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetUser gets user by id
// @Summary      Get an account
// @Description  This api is for getting user
// @Tags         user
// @Produce      json
// @Param        id   path      string  true  "Account ID"
// @Success      200  {object}  model.User
// @Router       /v1/users/{id} [get]
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
		h.log.Error("failed to get user", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteUser delete user by id
// @Summary      Delete an account
// @Description  This api is for delete user
// @Tags         user
// @Security BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Account ID"
// @Success      200  {object}  model.Id
// @Router       /v1/users/{id} [delete]
func (h *handlerV1) DeleteUser(c *gin.Context) {
	_, err := CheckClaims(h, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to auth token", logger.Error(err))
		return
	}
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	fmt.Println(guid)
	id, err := h.serviceManager.UserService().DeleteByID(
		ctx, &pb.GetIdFromUserID{Id: guid})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, id)
}
