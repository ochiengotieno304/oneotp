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
		Used:      false,
		ClientID:  clientID[0],
	})

	if err != nil {
		return nil, err
	}

	// encrypt data
	data := fmt.Sprintf("%s|%s|%s", fmt.Sprintf("%v", otpID.(primitive.ObjectID).Hex()), phone, otpCode) //ID|Phone|Code
	encrypted, err := utils.Encrypt(data)
	if err != nil {
		return nil, fmt.Errorf("error requesting otp: %v", err)
	}

	// send otp sms
	message := fmt.Sprintf("Your OTP code is %s, valid 15 minutes", otpCode)
	go sms.SendSMS(message, phone)

	return &pb.RequestOTPResponse{
		Ref: encrypted,
	}, nil
}

func (s *otpServer) VerifyOTP(ctx context.Context, req *pb.VerifyOTPRequest) (*pb.VerifyOTPResponse, error) {
	var reason string
	otpID := req.GetRef()
	var clientID []string = ctx.Value("clientID").([]string)

	decrypted, err := utils.Decrypt(otpID)
	if err != nil {
		return nil, fmt.Errorf("error verifying otp: %v", err)
	}

	info := strings.Split(decrypted, "|")
	otp, err := otpStore.FindOne(info[0], clientID[0])
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

	verifyExpired := time.Now().Before(otp.ExpiresAt)
	verifyPhone := info[1] == phone
	verifyCode := info[2] == code
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
		err = otpStore.UpdateOne(info[0], clientID[0])
		if err != nil {
			return nil, err
		}
	}

	return &pb.VerifyOTPResponse{
		Success: verifyExpired && verifyCode && verifyPhone && !verifyUsed,
		Reason:  reason,
	}, nil
}
