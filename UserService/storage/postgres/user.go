package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	pb "github.com/venomuz/service_api_swag_gin/UserService/genproto"
	"log"
)

type userRepo struct {
	db *sqlx.DB
}

//NewUserRepo ...
func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *pb.User) (*pb.User, error) {
	UserQuery := `INSERT INTO users(id,first_name,last_name,login,password,email,bio,phone_number,type_id,status) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := r.db.Exec(UserQuery, user.Id, user.FirstName, user.LastName, &user.Login, &user.Password, user.Email, user.Bio, user.PhoneNumber, user.TypeId, user.Status)
	if err != nil {
		log.Panicf("%s\n%s", "Error while users to table addresses", err)
	}
	fmt.Println("")
	AddressQuery := `INSERT INTO addresses(id,user_id,country,city,district,postal_code) VALUES($1,$2,$3,$4,$5,$6)`

	_, err = r.db.Exec(AddressQuery, user.Address.Id, user.Id, user.Address.Country, user.Address.City, user.Address.District, user.Address.PostalCode)
	if err != nil {
		log.Panicf("%s\n%s", "Error while inserting to table addresses", err)
	}

	return user, nil
}
func (r *userRepo) GetByID(ID string) (*pb.User, error) {
	user := pb.User{}
	GetUsers := `SELECT id, first_name, last_name, login, password, email, bio, phone_number, type_id, status FROM users WHERE id = $1`
	err := r.db.QueryRow(GetUsers, ID).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Login, &user.Password, &user.Email, &user.Bio, &user.PhoneNumber, &user.TypeId, &user.Status)
	if err != nil {
		return nil, err
	}

	addr := pb.Address{}
	GetAddresses := `SELECT id,user_id, city, district, country, postal_code FROM addresses WHERE user_id = $1`
	err = r.db.QueryRow(GetAddresses, user.Id).Scan(&addr.Id, &addr.UserId, &addr.City, &addr.District, &addr.Country, &addr.PostalCode)
	if err != nil {
		return nil, err
	}
	user.Address = &addr
	return &user, nil
}
func (r *userRepo) DeleteByID(ID string) (*pb.GetIdFromUserID, error) {
	id := pb.GetIdFromUserID{}
	err := r.db.QueryRow(`DELETE  FROM users WHERE id = $1 RETURNING id`, ID).Scan(&id.Id)
	if err != nil {
		log.Panicf("%s\n%s", "Error while deleteing data from table users", err)
	}

	return &id, nil
}
func (r *userRepo) GetAllUserFromDb(empty *pb.Empty) (*pb.AllUser, error) {
	var users pb.AllUser
	user := pb.User{}
	GetUsers := `SELECT id, first_name, last_name, login, password, email, bio, phone_number, type_id, status FROM users;`
	rows, err := r.db.Query(GetUsers)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Login, &user.Password, &user.Email, &user.Bio, &user.PhoneNumber, &user.TypeId, &user.Status)
		if err != nil {
			return nil, err
		}
		addr := pb.Address{}
		GetAddresses := `SELECT id,user_id, city, district, country, postal_code FROM addresses WHERE user_id = $1`
		err = r.db.QueryRow(GetAddresses, user.Id).Scan(&addr.Id, &addr.UserId, &addr.City, &addr.District, &addr.Country, &addr.PostalCode)
		if err != nil {
			return nil, err
		}
		user.Address = &addr
	}
	users.Users = append(users.Users, &user)

	return &users, nil
}
func (r *userRepo) GetList(page, limit int64) (*pb.LimitResponse, error) {
	offset := (page - 1) * limit
	fmt.Println(offset, page, limit)
	var userss pb.AllUser
	user := pb.User{}
	GetUsers := `SELECT id, first_name, last_name, login, password, email, bio, phone_number, type_id, status FROM users ORDER BY first_name OFFSET $1 LIMIT $2;`
	rows, err := r.db.Query(GetUsers, offset, limit)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Login, &user.Password, &user.Email, &user.Bio, &user.PhoneNumber, &user.TypeId, &user.Status)
		if err != nil {
			return nil, err
		}
		addr := pb.Address{}
		GetAddresses := `SELECT id,user_id, city, district, country, postal_code FROM addresses WHERE user_id = $1`
		err = r.db.QueryRow(GetAddresses, user.Id).Scan(&addr.Id, &addr.UserId, &addr.City, &addr.District, &addr.Country, &addr.PostalCode)
		if err != nil {
			return nil, err
		}
		user.Address = &addr
		userss.Users = append(userss.Users, &user)
	}

	var count int64
	CountUsersQuery := `SELECT count(*) FROM users`
	err = r.db.QueryRow(CountUsersQuery).Scan(&count)
	if err != nil {
		return nil, err
	}

	return &pb.LimitResponse{Users: userss.Users, AllUsers: count}, nil
}
func (r *userRepo) CheckValidLoginMail(key, value string) (bool, error) {
	c := 0
	if key == "login" {
		CheckQuery := `SELECT COUNT(1) FROM users WHERE login = $1`
		err := r.db.QueryRow(CheckQuery, value).Scan(&c)
		if c > 0 || err != nil {
			return true, err
		}
	} else if key == "email" {
		CheckQuery := `SELECT COUNT(1) FROM users WHERE email = $1`
		err := r.db.QueryRow(CheckQuery, value).Scan(&c)
		if c > 0 || err != nil {
			return true, err
		}
	}
	return false, nil
}
