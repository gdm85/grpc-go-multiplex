package main

/*
 * grpc-go-multiplex - proof-of-concept for gRPC/HTTP traffic multiplexing
 * Copyright (C) 2016 gdm85 - https://github.com/gdm85/grpc-go-multiplex/
This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
*/

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/cockroachdb/cmux"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

const port = 50051

func main() {
	// formatted address to listen on
	addr := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	// create the cmux object that will multiplex 2 protocols on same port
	m := cmux.New(l)
	// match gRPC requests, otherwise regular HTTP requests
	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := m.Match(cmux.Any())

	// create the go-grpc example greeter server
	grpcS := grpc.NewServer()
	pb.RegisterGreeterServer(grpcS, &server{})

	// create the regular HTTP requests muxer
	h := http.NewServeMux()
	h.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprintf(w, "OK")
    })
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            // The "/" pattern matches everything not matched by previous handlers
            fmt.Fprintf(w, "Welcome to the home page!")
    })
	httpS := &http.Server{
		Handler: h,
	}

	// collect on this channel the exits of each protocol's .Serve() call
	eps := make(chan error, 2)

	// start the listeners for each protocol
	go func() { eps <- grpcS.Serve(grpcL) }()
	go func() { eps <- httpS.Serve(httpL) }()

	log.Println("listening and serving (multiplexed) on", addr)
	err = m.Serve()

	// the rest of the code handles exit errors of the muxes

	var failed bool
	if err != nil {
		log.Println("cmux serve error: %v", err)
		failed = true
	}
	var i int
	for err := range eps {
		if err != nil {
			log.Printf("protocol serve error: %v", err)
			failed = true
		}
		i++
		if i == cap(eps) {
			close(eps)
			break
		}
	}
	if failed {
		os.Exit(1)
	}
}
