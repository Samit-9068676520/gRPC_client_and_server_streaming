package main

import (
	"context"
	"fmt"
	"greet/greetpb"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (res *greetpb.GreetRespnse, err error) {
	fmt.Printf("Greet function was invoked %v", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	result := "Hello" + firstName + lastName
	res = &greetpb.GreetRespnse{
		Result: result,
	}
	return res, nil
}
func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes was invoked %v", req)
	firstname := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello" + firstname + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}
func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("Long Greet function was invoke by the Client request:")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//we have finished reading the client stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("erroe while reading Client stream:  %v", err)
		}
		firstname := req.GetGreeting().GetFirstName()
		result += "Hello" + firstname + "! "
	}
}
func main() {
	fmt.Println("Hello I am Server:")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to Listem on %v :", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed To serve %v :", err)
	}
}
