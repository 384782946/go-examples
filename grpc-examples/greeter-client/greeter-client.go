package main

/*
note:关键代码是pb内已经根据协议生成了客户端，创建出对应客户端，并将预先生成的连接传递过去
*/

import (
	"log"
	"os"
	"time"

	pb "github.com/384782946/go-examples/grpc-examples/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defualtName = "world"
)

func main() {
	//生成一个连接
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	//
	name := defualtName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	ctx, cannel := context.WithTimeout(context.Background(), time.Second)
	defer cannel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("coundl not greet: %v", err)
	}

	log.Printf("Greeting %s", r.Message)

}
