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

type authServer struct {
	pb.UnimplementedAuthServiceServer
}

func NewAuthServer() *authServer {
	return &authServer{}
}

var authStore = stores.NewAuthStore()

func (s *authServer) RequestOTP(ctx context.Context, req *pb.RequestOTPRequest) (*pb.RequestOTPResponse, error) {
	phone := req.GetPhone()
	otp := strings.Join(utils.GenerateOTP(), "")
	after := time.Now().Add(15 * time.Minute)

	otpID, err := authStore.CreateOTP(&models.OTP{Code: otp, Phone: phone, ExpiresAt: after})
	if err != nil {
		return nil, err
	}

	message := fmt.Sprintf("Beba OTP is %s, expiry in 15 minutes", otp)

	go sms.SendSMS(message, phone)

	return &pb.RequestOTPResponse{
		Otp: &pb.OTP{
			Code: otp,
			Id:   fmt.Sprintf("%v", otpID.(primitive.ObjectID).Hex()),
		},
	}, nil

}

func (s *authServer) VerifyOTP(ctx context.Context, req *pb.VerifyOTPRequest) (*pb.VerifyOTPResponse, error) {
	phone := req.GetPhone()
	otpID := req.GetOtp().GetId()
	otpCode := req.GetOtp().GetCode()
	var reason string

	otp, err := authStore.FindOne(otpID)
	if err != nil {
		return nil, err
	}

	verifyExpired := time.Now().Before(otp.ExpiresAt)
	verifyCode := otpCode == otp.Code
	verifyPhone := phone == otp.Phone

	if !verifyExpired {
		reason = "verification code expired"
	} else if !verifyPhone || !verifyCode {
		reason = "invalid verification code provided"
	} else {
		reason = "none"
	}

	return &pb.VerifyOTPResponse{
		Success: verifyExpired && verifyCode && verifyPhone,
		Reason:  reason,
	}, nil
}
