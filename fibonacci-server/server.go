package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	pb "fibonacci-pe2/fibonacci"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	gRPCport = flag.Int("gRPC port", 50052, "The gRPC server port")
	restPort = flag.Int("REST port", 8081, "The REST API server port")
)

// gRPCserver is used to implement the fibonacci server
type gRPCserver struct {
	pb.UnimplementedFibonacciServer
}

// restServer is used to implement a rest server
type restServer struct{}

// restResponse contains the response object for the REST call
type restResponse struct {
	Total int
}

// CalculateFibonacciTotal is called from the gRPC server
func (s *gRPCserver) CalculateFibonacciTotal(ctx context.Context, in *pb.FibonacciRequest) (*pb.FibonacciResponse, error) {
	start := time.Now()
	log.Printf("Received calculate fibonacci total for: %d with rule: %s", in.GetMaximum(), in.GetRule())

	fibTotal := calcTotal(int(in.GetMaximum()), in.GetRule())

	elapsed := time.Since(start)
	log.Printf("Total calculated for: %d with rule: %s - Answer is %d - Total time taken is %s", in.GetMaximum(), in.GetRule(), fibTotal, elapsed)

	return &pb.FibonacciResponse{Total: uint64(fibTotal)}, nil
}

func (s *restServer) CalculateFibonacciTotal(w http.ResponseWriter, req *http.Request) {

	requestDecoder := json.NewDecoder(req.Body)
	var request pb.FibonacciRequest
	requestDecoder.Decode(&request)

	start := time.Now()

	resp := &restResponse{
		Total: calcTotal(int(request.Maximum), request.Rule),
	}

	elapsed := time.Since(start)

	log.Printf("Total calculated for: %d with rule: %s - Total time taken is %s", request.Maximum, request.Rule, elapsed)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func startGrpcServer(wg *sync.WaitGroup) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *gRPCport))
	if err != nil {
		log.Fatalf("GRPC port failed to list: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterFibonacciServer(s, &gRPCserver{})

	reflection.Register(s)

	log.Printf("grpc server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	wg.Done()
}

func startRestServer(wg *sync.WaitGroup) {
	s := &restServer{}

	mux := http.NewServeMux()
	mux.HandleFunc("/fibonacci", s.CalculateFibonacciTotal)

	log.Printf("rest server listening at %v", *restPort)

	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatal(err)
	}
}

func main() {

	flag.Parse()

	// Create a WaitGroup for both gRPC and REST server
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go startGrpcServer(wg)
	go startRestServer(wg)

	wg.Wait()
}

// calcTotal calculates the total of a fibonacci sequence up to a max, with a specified rule for each value
func calcTotal(max int, rule string) int {
	fibs := fibLoop(max)
	total := 0

	switch rule {
	case "even":
		for _, fib := range fibs {
			if fib%2 == 0 {
				total += fib
			}
		}
	case "odd":
		for _, fib := range fibs {
			if fib%2 != 0 {
				total += fib
			}
		}
	default:
		for _, fib := range fibs {
			total += fib
		}
	}

	return total

}

// fibLoop loops through all the fibonacci numbers in a sequence up to a maximum and appends them to an slice
func fibLoop(max int) []int {
	fibs := []int{1, 1}
	i1, i2, next := 1, 1, 0

	for {
		next = i1 + i2
		if next > max {
			break
		}

		i1 = i2
		i2 = next
		fibs = append(fibs, next)
	}

	return fibs
}
