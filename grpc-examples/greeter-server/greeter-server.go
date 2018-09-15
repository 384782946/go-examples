package main

/*
note:关键代码是声明一个符合pb服务定义的类型，实现相应代码
*/

import (
	"log"
	"net"
	"time"

	pb "github.com/384782946/go-examples/grpc-examples/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) SayHello(r *pb.HelloRequest, stream pb.Greeter_SayHelloServer) error {
	log.Print("request", r)
	t := time.NewTicker(time.Second * 3)
	for {
		<-t.C
		reply := pb.HelloReply{Message: "hello world, " + time.Now().Format("2006/01/02 15:04:05")}
		err := stream.Send(&reply)
		log.Print(reply)
		if err != nil {
			break
		}
	}
	t.Stop()
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
