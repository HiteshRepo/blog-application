package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/HiteshRepo/blog-application/global"
	"github.com/HiteshRepo/blog-application/proto"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"google.golang.org/grpc"
)

type authServer struct {
	userCollection *mongo.Collection
}

func Validations(in *proto.SignupRequest) error {

	username, email, password := in.GetUsername(), in.GetEmail(), in.GetPassword()

	emailRegex := regexp.MustCompile(global.EmailRegex)

	if len(username) < 4 || len(username) > 20 {
		return errors.New("Username should be greater that 4 and less than 20.")
	}
	if len(email) < 7 || len(email) > 35 {
		return errors.New("Email should be greater that 7 and less than 35.")
	}
	if len(password) < 8 || len(password) > 120 {
		return errors.New("Password should be greater that 8 and less than 120.")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("Invalid email format.")
	}
	return nil
}

func (a *authServer) Login(_ context.Context, in *proto.LoginRequest) (*proto.AuthResponse, error) {

	// fetch login and password from request
	login, password := in.GetLogin(), in.GetPassword()

	// fetch from db should not take more that 5 seconds
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()

	// look for user by creds entered
	var user global.User
	a.userCollection.FindOne(ctx, bson.M{"$or": []bson.M{bson.M{"username": login}, bson.M{"email": login}}}).Decode(&user)

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

	err := Validations(in)
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
	_, err = a.userCollection.InsertOne(ctx, newUser)
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
	a.userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)

	return &proto.UsedResponse{Used: user != global.NilUser}, nil
}

func (a *authServer) EmailUsed(_ context.Context, in *proto.EmailUsedRequest) (*proto.UsedResponse, error) {
	// fetch email from request
	email := in.GetEmail()

	// fetch from db should not take more that 5 seconds
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()

	var user global.User
	a.userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

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

	fmt.Println("Starting......")

	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(server, &authServer{userCollection: global.DB.Collection("user")})

	// GRPC listener ":5000"
	listener, err := net.Listen("tcp", "0.0.0.0:5000") //+os.Getenv("GRPCPORT"))
	if err != nil {
		fmt.Printf("Error creating listener : %v\n", err)
	}
	go func() {
		fmt.Println("GRPC server serving....")
		fmt.Printf("serving gRPC: %v\n", server.Serve(listener).Error())
	}()

	grpcWebServer := grpcweb.WrapServer(server)

	httpServer := &http.Server{
		// proxy port ":9001"
		Addr: "0.0.0.0:9001", //+ os.Getenv("PORT"),
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 {
				grpcWebServer.ServeHTTP(w, r)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")
				w.Header().Set("grpc-status", "")
				w.Header().Set("grpc-message", "")
				if grpcWebServer.IsGrpcWebRequest(r) {
					grpcWebServer.ServeHTTP(w, r)
				}
			}
		}), &http2.Server{}),
	}

	fmt.Println("Proxy server is going up....")
	fmt.Printf("serving proxy : %v\n", httpServer.ListenAndServe().Error())
}
