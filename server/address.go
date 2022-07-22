package main

import (
	"context"

	"gitlab.com/0xhyacinths/dscan/server/proto"
)

func (h *handler) TxByAddress(ctx context.Context, g *proto.TxByAddressRequest) (*proto.TxByAddressResponse, error) {
	txs, err := h.es.NormalTxByAddress(g.Address, nil, nil, 1, 25, true)
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
	return &proto.TxByAddressResponse{
		Txs: rv,
	}, nil
}
