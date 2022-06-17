package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	rr "github.com/gomodule/redigo/redis"
	_ "github.com/venomuz/service_api_swag_gin/ApiGateway/api/model"
	pb "github.com/venomuz/service_api_swag_gin/ApiGateway/genproto"
	l "github.com/venomuz/service_api_swag_gin/ApiGateway/pkg/logger"
	mail "github.com/venomuz/service_api_swag_gin/ApiGateway/pkg/mail"
	"google.golang.org/protobuf/encoding/protojson"
	"math/rand"
	"net/http"
	"strconv"
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
// @Summary      Get an account
// @Description  This api is for getting user
// @Tags         user
// @Produce      json
// @Param        id   path      string  true  "Account ID"
// @Success      200  {object}  model.Useri
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
		h.log.Error("failed to get user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteUser delete user by id
// @Summary      Delete an account
// @Description  This api is for delete user
// @Tags         user
// @Produce      json
// @Param        id   path      string  true  "Account ID"
// @Success      200  {object}  model.Id
// @Router       /v1/users/{id} [delete]
func (h *handlerV1) DeleteUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	fmt.Println(guid)
	_, err := h.serviceManager.UserService().DeleteByID(
		ctx, &pb.GetIdFromUserID{Id: guid})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}
	_ = h.redisStorage.Set("123qwe", "123")
	res, _ := h.redisStorage.Get("qwe")
	c.JSON(http.StatusOK, res)
}

// CheckReg CheckUserAnd creates users
// @Summary      Create an account with check
// @Description  This api is for creating user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param 		 user  body model.Useri true "user body"
// @Success      200  {string}  Ok
// @Router       /v1/users/check [post]
func (h *handlerV1) CheckReg(c *gin.Context) {
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

	response, err := h.serviceManager.UserService().CheckLoginMail(ctx, &pb.Check{Key: "login", Value: body.Login})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to check user", l.Error(err))
		return
	}
	num := rand.Intn(999999)
	str := strconv.Itoa(num)
	if response.Status == false {
		err := mail.SendMail(str, body.Email[0])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to send user", l.Error(err))
			return
		}
	}
	info, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to marshal", l.Error(err))
		return
	}

	err = h.redisStorage.SetWithTTL(str, string(info), 500)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to setting to redis user", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, response)
}

// Verify CheckReg Check and post
// @Summary      Create an account
// @Description  This api is for Create user
// @Tags         user
// @Produce      json
// @Param        code   path      string  true  "Verify Code"
// @Success      200  {object}  model.Code
// @Router       /v1/users/verify/{code} [post]
func (h *handlerV1) Verify(c *gin.Context) {
	body := pb.Useri{}
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true
	code := c.Param("code")
	res, err := h.redisStorage.Get(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get from redis user", l.Error(err))
		return
	}
	ff, err := rr.String(res, err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to change to string user", l.Error(err))
		return
	}
	err = json.Unmarshal([]byte(ff), &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to unmarshal", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	_, err = h.serviceManager.UserService().Create(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to posting to db", l.Error(err))
		return
	}
}
