package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mctofu/music-library-grpc/go/mlibgrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("failed with error: %v\n", err)
	}
}

func run() error {
	ctx := context.Background()

	conn, err := grpc.Dial("127.0.0.1:8337", grpc.WithInsecure())
	if err != nil {
		return err
	}
	client := mlibgrpc.NewMusicLibraryClient(conn)

	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	resp, err := client.Browse(ctx, &mlibgrpc.BrowseRequest{
		Path: path,
	})
	if err != nil {
		return err
	}

	dumpJSON(resp)

	return nil
}

func dumpJSON(v interface{}) error {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", string(jsonData))
	return nil
}
