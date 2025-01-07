package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/biosecret/daemon/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const socketPath = "/var/run/ztna.sock"

type DaemonServer struct {
	proto.UnimplementedDaemonServiceServer
}

func (s *DaemonServer) StartVPN(ctx context.Context, req *proto.Empty) (*proto.Response, error) {
	fmt.Println("Daemon nhận được lệnh StartVPN")
	return &proto.Response{Message: "VPN started"}, nil
}

func (s *DaemonServer) StopVPN(ctx context.Context, req *proto.Empty) (*proto.Response, error) {
	fmt.Println("Daemon nhận được lệnh StopVPN")
	return &proto.Response{Message: "VPN stopped"}, nil
}

func (s *DaemonServer) GetStatus(ctx context.Context, req *proto.Empty) (*proto.Response, error) {
	fmt.Println("Daemon nhận được lệnh GetStatus")
	return &proto.Response{Message: "VPN status: running"}, nil
}

func main() {
	// Xóa socket nếu đã tồn tại
	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Printf("Không thể tạo Unix socket: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()

	// Thay đổi quyền socket để user thông thường có thể truy cập
	os.Chmod(socketPath, 0660)

	server := grpc.NewServer()
	proto.RegisterDaemonServiceServer(server, &DaemonServer{})

	// Cho phép khám phá dịch vụ (hữu ích khi debug)
	reflection.Register(server)

	fmt.Println("Daemon đang chạy...")
	if err := server.Serve(listener); err != nil {
		fmt.Printf("Lỗi khi chạy server: %v\n", err)
	}
}
