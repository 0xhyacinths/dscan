package eshandler

import (
	"context"

	"github.com/nanmu42/etherscan-api"
	"gitlab.com/0xhyacinths/dscan/server/proto"
)

type EtherscanHandler struct {
	es *etherscan.Client
	proto.UnimplementedDescanIndexerServer
}

func NewEtherscanHandler(client *etherscan.Client) *EtherscanHandler {
	return &EtherscanHandler{
		es: client,
	}
}

func (h *EtherscanHandler) SayHi(ctx context.Context, g *proto.Hello) (*proto.Hello, error) {
	return &proto.Hello{
		Message: "Data provided by Etherscan.",
		ChainId: 1,
	}, nil
}
