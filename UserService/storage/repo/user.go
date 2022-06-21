package repo

import (
	pb "github.com/venomuz/service_api_swag_gin/UserService/genproto"
)

//UserStorageI ...
type UserStorageI interface {
	Create(user *pb.User) (*pb.User, error)
	GetByID(ID string) (*pb.User, error)
	DeleteByID(ID string) (*pb.GetIdFromUserID, error)
	GetAllUserFromDb(empty *pb.Empty) (*pb.AllUser, error)
	GetList(page, limit int64) (*pb.LimitResponse, error)
	CheckValidLoginMail(key, value string) (bool, error)
}
