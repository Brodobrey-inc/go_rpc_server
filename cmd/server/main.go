package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testGRPC/pkg/api"
	"testGRPC/pkg/cache"
	serverImpl "testGRPC/pkg/directory_informer"
	"time"

	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	c := cache.NewCache(time.Second*5, time.Second)

	wg.Add(1)
	go c.CleaningUp(ctx, &wg)

	s := grpc.NewServer()
	srv := serverImpl.NewDirectoryInformer(c)
	api.RegisterDirectoryInformerServer(s, srv)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		s.GracefulStop()
	}()

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}

	log.Println("Stopping server")
	cancel()
	wg.Wait()
	log.Println("Stopped")
}
