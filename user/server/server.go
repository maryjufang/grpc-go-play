package main

import (
	"context"
	"log"
	"net"

	"github.com/maryjufang/grpc-go-play/user/userpb"

	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

type UserManagementServer struct {
	conn                *pgx.Conn
	first_user_creation bool
	userpb.UnimplementedUserManagementServer
}

func (server *UserManagementServer) GetUsers(ctx context.Context, req *userpb.GetRequest) (*userpb.GetResponse, error) {
	createSql := `
	CREATE TABLE if not exists user_table(
		first_name varchar(80),
		last_name varchar(80),
		age int
	);
	`
	_, err := server.conn.Exec(context.Background(), createSql)
	if err != nil {
		log.Printf("ERROR in server.conn.Exec")
		return nil, err
	}

	var users_list *userpb.GetResponse = &userpb.GetResponse{}
	rows, err := server.conn.Query(context.Background(), "SELECT * FROM user_table;")
	if err != nil {
		log.Printf("ERROR in users query")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := userpb.User{}
		err = rows.Scan(&user.FirstName, &user.LastName, &user.Age)
		if err != nil {
			return nil, err
		}
		users_list.Users = append(users_list.Users, &user)

	}
	return users_list, nil
}

func (server *UserManagementServer) CreateUser(ctx context.Context, req *userpb.CreateRequest) (*userpb.CreateResponse, error) {
	createSql := `
	CREATE TABLE if not exists user_table(
		first_name varchar(80),
		last_name varchar(80),
		age int
	);
	`
	_, err := server.conn.Exec(context.Background(), createSql)
	if err != nil {
		log.Printf("ERROR in server.conn.Exec")
		return nil, err
	}

	server.first_user_creation = false

	log.Printf("Received: %v", req.User.GetFirstName())

	created_user := &userpb.User{FirstName: req.User.GetFirstName(), LastName: req.User.GetLastName(), Age: req.User.GetAge()}
	tx, err := server.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("conn.Begin failed: %v", err)
	}

	_, err = tx.Exec(context.Background(), "INSERT INTO user_table VALUES ($1,$2,$3)",
		created_user.FirstName, created_user.LastName, created_user.Age)
	if err != nil {
		log.Fatalf("tx.Exec failed: %v", err)
	}
	tx.Commit(context.Background())
	return &userpb.CreateResponse{User: created_user}, nil
}

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

func (server *UserManagementServer) Run() error {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	userpb.RegisterUserManagementServer(s, server)
	log.Printf("server listening at %v", lis.Addr())

	return s.Serve(lis)
}

func main() {
	log.Printf("In Users Manage Server func")
	database_url := "postgres://mary:marypassword@localhost:5432/marydb"
	var user_mgmt_server *UserManagementServer = NewUserManagementServer()

	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Fatalf("Unable to establish connection: %v", err)
	}
	defer conn.Close(context.Background())

	user_mgmt_server.conn = conn
	user_mgmt_server.first_user_creation = true
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
