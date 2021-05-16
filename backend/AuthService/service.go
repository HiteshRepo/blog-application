package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/HiteshRepo/blog-application/global"
	"github.com/HiteshRepo/blog-application/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"google.golang.org/grpc"
)

type authServer struct{}

func (a *authServer) Login(_ context.Context, in *proto.LoginRequest) (*proto.AuthResponse, error) {

	// fetch login and password from request
	login, password := in.GetLogin(), in.GetPassword()

	// fetch from db should not take more that 5 seconds
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()

	// look for user by creds entered
	var user global.User
	global.DB.Collection("user").FindOne(ctx, bson.M{"$or": []bson.M{bson.M{"username": login}, bson.M{"email": login}}}).Decode(&user)

	// check for empty user record
	if user == global.NilUser {
		return &proto.AuthResponse{}, errors.New("Invalid login credentials provided")
	}

	// validate password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return &proto.AuthResponse{}, errors.New("Invalid login credentials provided")
	}

	// send token
	return &proto.AuthResponse{Token: user.GetToken()}, nil
}

func (a *authServer) Signup(ctx context.Context, in *proto.SignupRequest) (*proto.AuthResponse, error) {

	err := a.Validations(in)
	if err != nil {
		return &proto.AuthResponse{}, errors.New(fmt.Sprintf("Validation failed : %s", err.Error()))
	}

	res, err := a.UsernameUsed(ctx, &proto.UsernameUsedRequest{Username: in.GetUsername()})
	if err != nil {
		log.Println("Error returned from UsernameUsed : ", err.Error())
		return nil, errors.New("Internal Error")
	}

	if res.GetUsed() {
		return nil, errors.New("Username already taken.")
	}

	res, err = a.EmailUsed(ctx, &proto.EmailUsedRequest{Email: in.GetEmail()})
	if err != nil {
		log.Println("Error returned from EmailUsed : ", err.Error())
		return nil, errors.New("Internal Error")
	}

	if res.GetUsed() {
		return nil, errors.New("Email already used.")
	}

	pw, _ := bcrypt.GenerateFromPassword([]byte(in.GetPassword()), bcrypt.DefaultCost)

	newUser := global.User{
		ID:       primitive.NewObjectID(),
		Email:    in.GetEmail(),
		Username: in.GetUsername(),
		Password: string(pw),
	}

	// insert user to db should not take more that 5 seconds
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	_, err = global.DB.Collection("user").InsertOne(ctx, newUser)
	if err != nil {
		log.Println("Error returned while inserting user to DB : ", err.Error())
		return nil, errors.New("Internal error while inserting user to DB.")
	}

	// send token
	return &proto.AuthResponse{Token: newUser.GetToken()}, nil
}

func (a *authServer) UsernameUsed(_ context.Context, in *proto.UsernameUsedRequest) (*proto.UsedResponse, error) {
	// fetch username from request
	username := in.GetUsername()

	// fetch from db should not take more that 5 seconds
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()

	var user global.User
	global.DB.Collection("user").FindOne(ctx, bson.M{"username": username}).Decode(&user)

	return &proto.UsedResponse{Used: user != global.NilUser}, nil
}

func (a *authServer) EmailUsed(_ context.Context, in *proto.EmailUsedRequest) (*proto.UsedResponse, error) {
	// fetch email from request
	email := in.GetEmail()

	// fetch from db should not take more that 5 seconds
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()

	var user global.User
	global.DB.Collection("user").FindOne(ctx, bson.M{"email": email}).Decode(&user)

	return &proto.UsedResponse{Used: user != global.NilUser}, nil
}

func (a *authServer) AuthUser(_ context.Context, in *proto.AuthUserRequest) (*proto.AuthUserResponse, error) {
	token := in.GetToken()
	user := global.UserFromToken(token)
	if user == global.NilUser {
		return &proto.AuthUserResponse{}, errors.New("Invalid token")
	}
	return &proto.AuthUserResponse{ID: user.ID.Hex(), Username: user.Username, Email: user.Email}, nil
}

func main() {
	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(server, &authServer{})

	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal("Error creating listener : ", err.Error())
	}
	server.Serve(listener)
}
