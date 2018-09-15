package main

/*
note:关键代码是pb内已经根据协议生成了客户端，创建出对应客户端，并将预先生成的连接传递过去
*/

import (
	"log"
	"os"
	"sync"

	pb "github.com/384782946/go-examples/grpc-examples/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defualtName = "world"
)

var wg sync.WaitGroup

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

	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("coundl not greet: %v", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			reply, _ := r.Recv()
			log.Print("receive reply: ", reply)
		}
	}()

	wg.Wait()
}
