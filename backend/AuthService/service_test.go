package main

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/HiteshRepo/blog-application/global"
	"github.com/HiteshRepo/blog-application/proto"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

func setup() {
	// connects to test-db
	global.ConnectToTestDatabase()
	// creates a user collection
	userCollection = global.DB.Collection("user")
}

func Test_authServer_Login(t *testing.T) {

	// test-password
	pw, _ := bcrypt.GenerateFromPassword([]byte("test-password"), bcrypt.DefaultCost)

	// insert to verify Login rpc functionality
	_, err := userCollection.InsertOne(context.Background(), global.User{
		ID:       primitive.NewObjectID(),
		Email:    "test-user@gmail.com",
		Username: "test-user",
		Password: string(pw),
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	testcases := []map[string]interface{}{
		map[string]interface{}{
			"login":    "test-user@gmail.com",
			"password": "test-password",
		},
		map[string]interface{}{
			"login":    "incorrect-email@gmail.com",
			"password": "incorrect-password",
			"error":    "Invalid login credentials provided",
		},
		map[string]interface{}{
			"login":    "test-user",
			"password": "test-password",
		},
		map[string]interface{}{
			"login":    "incorrect-username",
			"password": "incorrect-password",
			"error":    "Invalid login credentials provided",
		},
	}

	server := authServer{}

	for _, tcase := range testcases {

		resp, err := server.Login(context.Background(), &proto.LoginRequest{Login: tcase["login"].(string), Password: tcase["password"].(string)})

		if errMsg, ok := tcase["error"]; ok {
			// invalid creds
			assert.Errorf(t, err, "case: %v", tcase)
			assert.Containsf(t, err.Error(), errMsg.(string), "case: %v", tcase)
		} else {
			// valid creds - email/username & password
			assert.NoError(t, err, "case: %v", tcase)
			assert.Truef(t, len(resp.GetToken()) > 0, "case: %v", tcase)
		}
	}
}

func Test_authServer_UsernameUsed(t *testing.T) {

	// insert to verify UsernameUsed rpc functionality
	_, err := userCollection.InsertOne(context.Background(), global.User{ID: primitive.NewObjectID(), Username: "Test-UserName-Used"})
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	testCases := []map[string]interface{}{
		map[string]interface{}{
			"username": "Test-UserName-Used",
			"used":     true,
		},
		map[string]interface{}{
			"username": "Test-UserName-Unused",
			"used":     false,
		},
	}

	server := authServer{}

	for _, tcase := range testCases {

		res, err := server.UsernameUsed(context.Background(), &proto.UsernameUsedRequest{Username: tcase["username"].(string)})
		if !assert.NoErrorf(t, err, "case: %v", tcase) {
			t.FailNow()
		}
		assert.Truef(t, res.GetUsed() == tcase["used"].(bool), "case: %v", tcase)
	}
}

func Test_authServer_EmailUsed(t *testing.T) {
	// insert to verify EmailUsed rpc functionality
	_, err := userCollection.InsertOne(context.Background(), global.User{ID: primitive.NewObjectID(), Email: "Test-Email-Used@gmail.com"})
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	testCases := []map[string]interface{}{
		map[string]interface{}{
			"email": "Test-Email-Used@gmail.com",
			"used":  true,
		},
		map[string]interface{}{
			"email": "Test-Email-Unused@gmail.com",
			"used":  false,
		},
	}

	server := authServer{}

	for _, tcase := range testCases {

		res, err := server.EmailUsed(context.Background(), &proto.EmailUsedRequest{Email: tcase["email"].(string)})
		if !assert.NoErrorf(t, err, "case: %v", tcase) {
			t.FailNow()
		}
		assert.Truef(t, res.GetUsed() == tcase["used"].(bool), "case: %v", tcase)
	}
}

func Test_authServer_Signup(t *testing.T) {
	// test-password
	pw, _ := bcrypt.GenerateFromPassword([]byte("test-signup-password"), bcrypt.DefaultCost)

	// insert to verify Login rpc functionality
	_, err := userCollection.InsertOne(context.Background(), global.User{
		ID:       primitive.NewObjectID(),
		Email:    "test-signup-user@gmail.com",
		Username: "test-signup-user",
		Password: string(pw),
	})
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	testCases := []map[string]interface{}{
		map[string]interface{}{
			"username": "test-signup-user",
			"email":    "test-signup-user2@gmail.com",
			"password": "test-signup-password",
			"error":    "Username already taken.",
		},
		map[string]interface{}{
			"username": "test-signup-user2",
			"email":    "test-signup-user@gmail.com",
			"password": "test-signup-password",
			"error":    "Email already used.",
		},
		map[string]interface{}{
			"username": "test-signup-user2",
			"email":    "test-signup-user2@gmail.com",
			"password": "test-signup-password",
		},
	}

	server := authServer{}

	for _, tcase := range testCases {

		resp, err := server.Signup(context.Background(), &proto.SignupRequest{
			Username: tcase["username"].(string),
			Email:    tcase["email"].(string),
			Password: tcase["password"].(string),
		})

		if errMsg, ok := tcase["error"]; ok {
			assert.Errorf(t, err, "case: %v", tcase)
			assert.Containsf(t, err.Error(), errMsg, "case: %v", tcase)
		} else {
			assert.NoErrorf(t, err, "case: %v", tcase)
			assert.Truef(t, len(resp.GetToken()) > 0, "case: %v", tcase)
		}
	}
}

func Test_authServer_AuthUser(t *testing.T) {

	server := authServer{}
	token := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoie1wiSURcIjpcIjYwYTBjNDA3ZWQyNmFhNjViNmI2MDlmZVwiLFwiVXNlcm5hbWVcIjpcInRlc3QtdXNlclwiLFwiRW1haWxcIjpcInRlc3QtdXNlckBnbWFpbC5jb21cIixcIlBhc3N3b3JkXCI6XCIkMmEkMTAkSFMuZzFCTGROLnU5d1JXcUpReU53dVVIN3N2Q3lULkdramVLUmNqTVA5S0hQSUdDbG9qVy5cIn0ifQ.9kCEVYGfbv4ra_5tPpzCT9AT9YQEfn4AKsDNUlV33iw`
	email := "test-user@gmail.com"
	username := "test-user"

	testCases := []map[string]interface{}{
		map[string]interface{}{
			"token": token,
		},
		map[string]interface{}{
			"token": "incorrect-auth-token",
			"error": "Invalid token",
		},
	}

	for _, tcase := range testCases {

		resp, err := server.AuthUser(context.Background(), &proto.AuthUserRequest{Token: tcase["token"].(string)})

		if errMsg, ok := tcase["error"]; ok {
			assert.Errorf(t, err, "case: %v", tcase)
			assert.Containsf(t, err.Error(), errMsg.(string), "case: %v", tcase)
		} else {
			assert.NoErrorf(t, err, "case: %v", tcase)
			assert.Truef(t, len(resp.GetID()) > 0, "case: %v", tcase)
			assert.Truef(t, resp.GetEmail() == email, "case: %v", tcase)
			assert.Truef(t, resp.GetUsername() == username, "case: %v", tcase)
		}
	}
}

func teardown() {

	// deleting all records in a cluster but not the cluster
	// result, err := userCollection.DeleteMany(context.Background(), bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("DeleteMany removed %v document(s)\n", result.DeletedCount)

	//dropping user collection after testing
	if err := userCollection.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {

	// defer function to invoke teardown after catching any panic
	defer func() {
		if panicErr := recover(); panicErr != nil {
			log.Fatal(panicErr)
			teardown()
		}
	}()

	setup()

	retCode := m.Run()

	teardown()

	os.Exit(retCode)
}
