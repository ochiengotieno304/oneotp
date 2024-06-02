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

var authStore = stores.NewAuthStore()

func (s *otpServer) RequestOTP(ctx context.Context, req *pb.RequestOTPRequest) (*pb.RequestOTPResponse, error) {
	phone := req.GetPhone()
	otp := strings.Join(utils.GenerateOTP(), "")
	after := time.Now().Add(15 * time.Minute)

	otpID, err := authStore.CreateOTP(&models.OTP{Code: otp, Phone: phone, ExpiresAt: after, Used: false})
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

func (s *otpServer) VerifyOTP(ctx context.Context, req *pb.VerifyOTPRequest) (*pb.VerifyOTPResponse, error) {
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
	verifyUsed := otp.Used

	if !verifyExpired {
		reason = "verification code expired"
	} else if verifyUsed {
		reason = "code already used"
	} else if !verifyPhone || !verifyCode {
		reason = "invalid verification code provided"
	} else {
		reason = "none"
	}

	if !verifyUsed {
		err = authStore.UpdateOne(otpID)
		if err != nil {
			return nil, err
		}
	}

	return &pb.VerifyOTPResponse{
		Success: verifyExpired && verifyCode && verifyPhone && !verifyUsed,
		Reason:  reason,
	}, nil
}
