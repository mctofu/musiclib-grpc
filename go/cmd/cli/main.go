package main

import (
	"encoding/json"
	"errors"
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

	if len(os.Args) < 2 {
		return errors.New("must specify command")
	}

	cmd := os.Args[1]

	switch cmd {
	case "browse":
		return browse(ctx, client, os.Args[2:])
	case "media":
		return media(ctx, client, os.Args[2:])
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}

}

func browse(ctx context.Context, client mlibgrpc.MusicLibraryClient, args []string) error {
	var path string

	if len(args) > 0 {
		path = args[0]
	}

	resp, err := client.Browse(ctx, &mlibgrpc.BrowseRequest{
		Uri: path,
	})
	if err != nil {
		return err
	}

	dumpJSON(resp)

	return nil
}

func media(ctx context.Context, client mlibgrpc.MusicLibraryClient, args []string) error {
	var path string

	if len(args) > 0 {
		path = args[0]
	}

	resp, err := client.Media(ctx, &mlibgrpc.MediaRequest{
		Uri: path,
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
