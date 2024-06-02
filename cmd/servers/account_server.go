package servers

import (
	"context"

	"github.com/ochiengotieno304/oneotp/pkg/db/models"
	"github.com/ochiengotieno304/oneotp/pkg/db/stores"
	"github.com/ochiengotieno304/oneotp/pkg/pb"
)

type accountServer struct {
	pb.UnimplementedAccountServiceServer
}

func NewAccountServer() *accountServer {
	return &accountServer{}
}

var accountStore = stores.NewAccountStore()

func (s *accountServer) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	name := req.GetAccount().GetName()
	phone := req.GetAccount().GetPhone()

	err := accountStore.CreateAccount(&models.Account{
		Name:  name,
		Phone: phone,
	})

	if err != nil {
		return nil, err
	}

	return &pb.CreateAccountResponse{
		Account: &pb.Account{
			Name:  name,
			Phone: phone,
		},
	}, nil
}
