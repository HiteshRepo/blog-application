package main

import (
	"errors"
	"regexp"

	"github.com/HiteshRepo/blog-application/global"
	"github.com/HiteshRepo/blog-application/proto"
)

func (a *authServer) Validations(in *proto.SignupRequest) error {

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
