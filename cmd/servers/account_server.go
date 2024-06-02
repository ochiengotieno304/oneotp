package servers

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/ochiengotieno304/oneotp/internal/helpers/auth"
	"github.com/ochiengotieno304/oneotp/internal/helpers/email"
	"github.com/ochiengotieno304/oneotp/pkg/db/models"
	"github.com/ochiengotieno304/oneotp/pkg/db/stores"
	"github.com/ochiengotieno304/oneotp/pkg/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type accountServer struct {
	pb.UnimplementedAccountServiceServer
}

func NewAccountServer() *accountServer {
	return &accountServer{}
}

var accountStore = stores.NewAccountStore()

func validatePassword(password string, altPassword string) bool {
	return password == altPassword
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (s *accountServer) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	name := req.GetAccount().GetName()
	phone := req.GetAccount().GetPhone()
	email := req.GetAccount().GetEmail()
	password := req.GetAccount().GetPassword()
	altPassword := req.GetAccount().GetAltPassword()

	if name == "" {
		return nil, fmt.Errorf("name cannot be blank")
	} else if email == "" {
		return nil, fmt.Errorf("email cannot be blank")
	} else if password == "" || altPassword == "" {
		return nil, fmt.Errorf("password cannot be blank")
	}

	if !validateEmail(email) {
		return nil, fmt.Errorf("invalid email address")
	}

	if !validatePassword(password, altPassword) {
		return nil, fmt.Errorf("mismatch in passwords")
	}

	hashedPass, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	accountID, err := accountStore.CreateAccount(&models.Account{
		Name:         name,
		Phone:        phone,
		Email:        email,
		PasswordHash: hashedPass,
	})
	if err != nil {
		return nil, err
	}

	accessToken, err := auth.GenerateToken(accountID.(primitive.ObjectID).Hex())
	if err != nil {
		return nil, err
	}

	go mail_helper.SendMail([]string{email}, []byte(""), []byte("Welcome To OneOTP"))

	return &pb.CreateAccountResponse{
		AccessToken: accessToken,
	}, nil
}
