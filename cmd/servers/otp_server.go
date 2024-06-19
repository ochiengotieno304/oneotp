package servers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ochiengotieno304/oneotp/internal/helpers/sms"
	"github.com/ochiengotieno304/oneotp/internal/utils"
	"github.com/ochiengotieno304/oneotp/pkg/db/models"
	"github.com/ochiengotieno304/oneotp/pkg/db/stores"
	"github.com/ochiengotieno304/oneotp/pkg/pb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type otpServer struct {
	pb.UnimplementedOTPServiceServer
}

func NewOTPServer() *otpServer {
	return &otpServer{}
}

var otpStore = stores.NewOTPStore()

func (s *otpServer) RequestOTP(ctx context.Context, req *pb.RequestOTPRequest) (*pb.RequestOTPResponse, error) {
	phone := req.GetPhone()
	otpCode := strings.Join(utils.GenerateOTP(), "")
	after := time.Now().Add(15 * time.Minute)
	var clientID []string = ctx.Value("clientID").([]string)

	// encrypt data before storage
	code, err := utils.Encrypt(otpCode)
	if err != nil {
		return nil, err
	}
	phoneEncrpted, err := utils.Encrypt(phone)
	if err != nil {
		return nil, err
	}

	otpID, err := otpStore.CreateOTP(&models.OTP{
		Code:      code,
		Phone:     phoneEncrpted,
		ExpiresAt: after,
		Attempts:  0,
		ClientID:  clientID[0],
	})

	if err != nil {
		return nil, err
	}

	// send otp sms
	message := fmt.Sprintf("Your OTP code is %s, valid 15 minutes", otpCode)
	go sms.SendSMS(message, phone)

	return &pb.RequestOTPResponse{
		Ref: otpID.(primitive.ObjectID).Hex(),
	}, nil
}

func (s *otpServer) VerifyOTP(ctx context.Context, req *pb.VerifyOTPRequest) (*pb.VerifyOTPResponse, error) {
	var reason string
	otpID := req.GetRef()
	codeFromRequest := req.GetCode()
	phoneFromRequest := req.GetPhone()

	var clientID []string = ctx.Value("clientID").([]string)

	otp, err := otpStore.FindOne(otpID, clientID[0])
	if err != nil {
		return nil, err
	}

	code, err := utils.Decrypt(otp.Code)
	if err != nil {
		return nil, err
	}

	phone, err := utils.Decrypt(otp.Phone)
	if err != nil {
		return nil, err
	}

	expired := time.Now().After(otp.ExpiresAt) || otp.Expired
	notPhone := phoneFromRequest != phone
	notCode := codeFromRequest != code
	exceededAttempts := otp.Attempts == 3

	if expired {
		reason = "verification code expired"
	} else if exceededAttempts {
		reason = "max attempts reached, please request a new otp"
	} else if notPhone && notCode {
		reason = "invalid verification code provided"
	} else {
		reason = "none"
	}

	if !exceededAttempts {
		err = otpStore.UpdateOne(otpID, clientID[0], 1) // increament attempsts
		if err != nil {
			return nil, err
		}
	}

	if !expired && !notCode && !notPhone && !exceededAttempts {
		err = otpStore.UpdateOne(otpID, clientID[0], 2) // expire otp after verifiction
		if err != nil {
			return nil, err
		}
	}

	return &pb.VerifyOTPResponse{
		Success: !expired && !notCode && !notPhone && !exceededAttempts,
		Reason:  reason,
	}, nil
}
