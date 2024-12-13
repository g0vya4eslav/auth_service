package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "auth-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedAuthServiceServer
}

func (s *server) Authenticate(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	if req.Username == "admin" && req.Password == "password" {
		return &pb.AuthResponse{Success: true}, nil
	}
	return &pb.AuthResponse{Success: false}, fmt.Errorf("неверные учетные данные")
}

func authRequestHandler(w http.ResponseWriter, r *http.Request) {
	authUser := r.Header.Get("Auth-User")
	authPass := r.Header.Get("Auth-Pass")

	log.Printf("Authenticating user: %s", authUser)

	// Проверка учетных данных
	if authUser == "" || authPass == "" || !(authUser == "admin" && authPass == "password") {
		w.Header().Set("Auth-Status", "Invalid login or password")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Auth-Status", "OK")
	w.Header().Set("Auth-Server", "127.0.0.1") // Сервер IMAP
	w.Header().Set("Auth-Port", "1993")        // Порт IMAP
	w.WriteHeader(http.StatusOK)
}

func main() {
	go func() {
		listener, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Не удалось запустить gRPC-сервер: %v", err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterAuthServiceServer(grpcServer, &server{})
		reflection.Register(grpcServer)

		fmt.Println("gRPC сервер запущен на порту 50051")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Ошибка запуска gRPC-сервера: %v", err)
		}
	}()

	http.HandleFunc("/auth", authRequestHandler)
	fmt.Println("HTTP сервер запущен на порту 9000 для auth_request")
	log.Fatal(http.ListenAndServe("0.0.0.0:9000", nil))
}
