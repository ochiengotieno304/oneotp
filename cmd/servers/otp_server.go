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
	otp := strings.Join(utils.GenerateOTP(), "")
	after := time.Now().Add(15 * time.Minute)
	var clientID []string = ctx.Value("clientID").([]string)

	otpID, err := otpStore.CreateOTP(&models.OTP{
		Code:      utils.Encrypt(otp),
		Phone:     utils.Encrypt(phone),
		ExpiresAt: after,
		Used:      false,
		ClientID:  clientID[0],
	})

	if err != nil {
		return nil, err
	}

	message := fmt.Sprintf("Your OTP code is %s, valid 15 minutes", otp) // send verification sms

	go sms.SendSMS(message, phone)

	return &pb.RequestOTPResponse{
		Id: fmt.Sprintf("%v", otpID.(primitive.ObjectID).Hex()),
	}, nil

}

func (s *otpServer) VerifyOTP(ctx context.Context, req *pb.VerifyOTPRequest) (*pb.VerifyOTPResponse, error) {
	phone := req.GetPhone()
	otpID := req.GetId()
	otpCode := req.GetCode()
	var reason string

	otp, err := otpStore.FindOne(otpID)
	if err != nil {
		return nil, err
	}

	verifyExpired := time.Now().Before(otp.ExpiresAt)
	verifyCode := otpCode == utils.Decrypt(otp.Code)
	verifyPhone := phone == utils.Decrypt(otp.Phone)
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
		err = otpStore.UpdateOne(otpID)
		if err != nil {
			return nil, err
		}
	}

	return &pb.VerifyOTPResponse{
		Success: verifyExpired && verifyCode && verifyPhone && !verifyUsed,
		Reason:  reason,
	}, nil
}
