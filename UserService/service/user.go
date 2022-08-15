package service

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	pb "github.com/venomuz/service-api-swag-gin/UserService/genproto"
	l "github.com/venomuz/service-api-swag-gin/UserService/pkg/logger"
	cl "github.com/venomuz/service-api-swag-gin/UserService/service/grpc_client"
	"github.com/venomuz/service-api-swag-gin/UserService/storage"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//UserService ...
type UserService struct {
	storage storage.IStorage
	logger  l.Logger
	client  cl.GrpcClientI
}

//NewUserService ...
func NewUserService(db *sqlx.DB, log l.Logger, client cl.GrpcClientI) *UserService {
	return &UserService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
	}
}

func HashPassword(password string) (string, error) {
	pw := []byte(password)
	result, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (s *UserService) Create(ctx context.Context, req *pb.User) (*pb.User, error) {
	req.Id = uuid.NewV4().String()
	req.Address.Id = uuid.NewV4().String()
	passReq, err := HashPassword(req.Password)
	if err != nil {
		s.logger.Error("Error while Encrypt user password", l.Error(err))
		return nil, status.Error(codes.Internal, "Error insert user")
	}
	req.Password = passReq
	user, err := s.storage.User().Create(req)
	if err != nil {
		s.logger.Error("Error while inserting user info", l.Error(err))
		return nil, status.Error(codes.Internal, "Error insert user")
	}

	for _, posts := range user.Posts {
		posts.UserId = req.Id
		post, err := s.client.PostService().PostCreate(ctx, posts)
		if err != nil {
			return nil, err
		}
		user.Posts = append(user.Posts, post)
	}

	return user, err
}
func (s *UserService) GetByID(ctx context.Context, req *pb.GetIdFromUser) (*pb.User, error) {
	user, err := s.storage.User().GetByID(req.Id)
	if err != nil {
		fmt.Println(err)
		s.logger.Error("Error while getting user info", l.Error(err))
		return nil, status.Error(codes.Internal, "Error insert user")
	}

	post, err := s.client.PostService().PostGetAllPosts(ctx, req)
	user.Posts = post.Posts
	return user, err
}
func (s *UserService) DeleteByID(ctx context.Context, req *pb.GetIdFromUserID) (*pb.GetIdFromUserID, error) {

	user, err := s.storage.User().DeleteByID(req.Id)
	if err != nil {
		s.logger.Error("Error while getting user info", l.Error(err))
		return nil, status.Error(codes.Internal, "Error insert user")
	}

	_, err = s.client.PostService().PostDeleteByID(ctx, &pb.GetIdFromUser{Id: user.Id})
	return user, err
}
func (s *UserService) GetAllByUserId(ctx context.Context, req *pb.GetIdFromUser) (*pb.Post, error) {
	post, err := s.client.PostService().PostGetByID(ctx, req)
	if err != nil {
		fmt.Println(err)
		s.logger.Error("Error while getting post info", l.Error(err))
		return nil, status.Error(codes.Internal, "Error insert post")
	}

	return post, err
}
func (s *UserService) GetAllUserFromDb(ctx context.Context, req *pb.Empty) (*pb.AllUser, error) {
	users, err := s.storage.User().GetAllUserFromDb(req)
	if err != nil {
		fmt.Println(err)
		s.logger.Error("Error while getting post info", l.Error(err))
		return nil, status.Error(codes.Internal, "Error insert post")
	}

	user := users.Users
	for _, usr := range user {
		aa := pb.GetIdFromUser{}
		aa.Id = usr.Id
		post, err := s.client.PostService().PostGetAllPosts(ctx, &aa)
		if err != nil {
			fmt.Println(err)
			s.logger.Error("Error while getting post info", l.Error(err))
			return nil, status.Error(codes.Internal, "Error insert post")
		}
		usr.Posts = post.Posts

	}
	users.Users = user

	return users, err
}
func (s *UserService) GetList(ctx context.Context, req *pb.LimitRequest) (*pb.LimitResponse, error) {
	users, err := s.storage.User().GetList(req.Page, req.Limit)
	if err != nil {
		fmt.Println(err)
		s.logger.Error("Error while getting post info", l.Error(err))
		return nil, status.Error(codes.Internal, "Error insert post")
	}

	user := users.Users
	for _, usr := range user {
		aa := pb.GetIdFromUser{}
		aa.Id = usr.Id
		post, err := s.client.PostService().PostGetAllPosts(ctx, &aa)
		if err != nil {
			fmt.Println(err)
			s.logger.Error("Error while getting post info", l.Error(err))
			return nil, status.Error(codes.Internal, "Error insert post")
		}
		usr.Posts = post.Posts

	}
	users.Users = user

	return users, err
}
func (s *UserService) CheckLoginMail(_ context.Context, check *pb.Check) (*pb.Okay, error) {
	get, err := s.storage.User().CheckValidLoginMail(check.Key, check.Value)
	if err != nil {
		fmt.Println(err)
		s.logger.Error("Error while CheckLoginMail user info", l.Error(err))
		return nil, status.Error(codes.Internal, "Error insert v")
	}
	return &pb.Okay{Status: get}, err
}

func (s *UserService) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	res, id, err := s.storage.User().Login(request.Mail, request.Password)
	if err != nil {
		fmt.Println(err)
		s.logger.Error("Error while get user via mail from db", l.Error(err))
		return nil, status.Error(codes.Internal, "Error insert v")
	}

	user, err := s.storage.User().GetByID(id)
	if err != nil {
		fmt.Println(err)
		s.logger.Error("Error while get user via id from db", l.Error(err))
		return nil, status.Error(codes.Internal, "Error insert v")
	}
	return &pb.LoginResponse{Check: res, UserData: user}, nil
}
