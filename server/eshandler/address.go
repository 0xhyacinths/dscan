package eshandler

import (
	"context"

	"github.com/0xhyacinths/dscan/server/proto"
)

func (h *EtherscanHandler) TxByAddress(ctx context.Context, g *proto.TxByAddressRequest) (*proto.TxByAddressResponse, error) {
	txs, err := h.es.NormalTxByAddress(g.Address, nil, nil, int(g.Page), int(g.Offset)+1, true)
	if err != nil {
		return nil, err
	}

	rv := make([]*proto.TransactionSummary, len(txs))
	for i, tx := range txs {
		rv[i] = &proto.TransactionSummary{
			From:   tx.From,
			To:     tx.To,
			Amount: tx.Value.Int().String(),
			Tx:     tx.Hash,
			Block:  uint64(tx.BlockNumber),
		}
	}

	hasMore := false
	if len(rv) > int(g.Offset) {
		rv = rv[:len(rv)-1]
		hasMore = true
	}

	return &proto.TxByAddressResponse{
		Txs:     rv,
		HasMore: hasMore,
	}, nil
}
