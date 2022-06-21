package v1

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	_ "github.com/venomuz/service_api_swag_gin/ApiGateway/api/model"
	pb "github.com/venomuz/service_api_swag_gin/ApiGateway/genproto"
	l "github.com/venomuz/service_api_swag_gin/ApiGateway/pkg/logger"
	"github.com/venomuz/service_api_swag_gin/ApiGateway/pkg/mail"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"google.golang.org/protobuf/encoding/protojson"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// CheckReg CheckUserAnd creates users
// @Summary      Create an account with check
// @Description  This api is for creating user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param 		 user  body model.User true "user body"
// @Success      200  {string}  Ok
// @Router       /v1/users/check [post]
func (h *handlerV1) CheckReg(c *gin.Context) {
	var (
		JsMarshal protojson.MarshalOptions
		body      pb.User
	)

	JsMarshal.UseProtoNames = true

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

	const minEntropyBits = 60
	err = passwordvalidator.Validate(body.Password, minEntropyBits)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to valid password user", l.Error(err))
		return
	}
	
	num := rand.Intn(999999)
	str := strconv.Itoa(num)
	if response.Status == false {
		err := mail.SendMail(str, body.Email)
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
	body := pb.User{}
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
	ff, err := redis.String(res, err)
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
