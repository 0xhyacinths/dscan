package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nanmu42/etherscan-api"
	"gitlab.com/0xhyacinths/dscan/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
)

func main() {
	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	hd := &handler{
		es: etherscan.New(etherscan.Mainnet, os.Getenv("ES_API")),
	}
	proto.RegisterDescanIndexerServer(grpcServer, hd)
	go func() {
		log.Fatalln(grpcServer.Serve(lis))
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:9091",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = proto.RegisterDescanIndexerHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":9090",
		Handler: corsWrapper(gwmux),
	}

	fmt.Println("starting")
	if err := gwServer.ListenAndServe(); err != nil {
		grpclog.Fatalf("failed starting http server: %v", err)
	}
}

func corsWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

type handler struct {
	es *etherscan.Client
	proto.UnimplementedDescanIndexerServer
}

func (h *handler) SayHi(ctx context.Context, g *proto.Hello) (*proto.Hello, error) {
	return &proto.Hello{
		Message: "Data provided by Etherscan.",
	}, nil
}
