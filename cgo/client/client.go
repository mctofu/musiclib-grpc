package main

/*
typedef struct MLibGRPC_BrowseItem {
	char *name;
	char *uri;
	int folder;
} MLibGRPC_BrowseItem;
*/
import "C"

import (
	"context"
	"errors"
	"log"
	"unsafe"

	"github.com/mctofu/music-library-grpc/go/mlibgrpc"
	"google.golang.org/grpc"
)

var client mlibgrpc.MusicLibraryClient
var conn *grpc.ClientConn

//export MLibGRPC_Connect
func MLibGRPC_Connect() int {
	if client != nil {
		log.Printf("MLibGRPC_Connect: already connected\n")
		return 1
	}

	grpcConn, err := grpc.Dial("127.0.0.1:8337", grpc.WithInsecure())
	if err != nil {
		log.Printf("MLibGRPC_Connect: failed to connect: %v\n", err)
		return 2
	}
	conn = grpcConn
	client = mlibgrpc.NewMusicLibraryClient(conn)

	log.Printf("MLibGRPC_Connect: Connected!")

	return 0
}

//export MLibGRPC_Disconnect
func MLibGRPC_Disconnect() int {
	if client == nil {
		log.Printf("MLibGRPC_Disconnect: not connected\n")
		return 1
	}

	if err := conn.Close(); err != nil {
		log.Printf("MLibGRPC_Disconnect: failed to close: %v\n", err)
		return 2
	}

	client = nil
	conn = nil

	log.Printf("Disconnected!")

	return 0
}

//export MLibGRPC_Browse
func MLibGRPC_Browse(cPath *C.char, cSearch *C.char) **C.struct_MLibGRPC_BrowseItem {
	path := C.GoString(cPath)
	search := C.GoString(cSearch)

	log.Printf("MLibGRPC_Browse: Browse path: %s search: %s\n", path, search)

	items, err := browse(path, search)
	if err != nil {
		log.Printf("MLibGRPC_Browse: error %v\n", err)
	}

	result := C.malloc(C.size_t(len(items)+1) * C.size_t(unsafe.Sizeof(uintptr(0))))
	resultArr := (*[1<<30 - 1]*C.struct_MLibGRPC_BrowseItem)(result)

	for i, item := range items {
		cBrowseItem := (*C.struct_MLibGRPC_BrowseItem)(C.malloc(
			C.size_t(unsafe.Sizeof(C.struct_MLibGRPC_BrowseItem{}))))
		cBrowseItem.name = C.CString(item.Name)
		cBrowseItem.uri = C.CString(item.Uri)
		if item.Folder {
			cBrowseItem.folder = 1
		} else {
			cBrowseItem.folder = 0
		}

		resultArr[i] = cBrowseItem
	}

	resultArr[len(items)] = nil

	log.Printf("MLibGRPC_Browse: returning\n")

	return (**C.struct_MLibGRPC_BrowseItem)(result)
}

func browse(path string, search string) ([]*mlibgrpc.BrowseItem, error) {
	if client == nil {
		log.Printf("MLibGRPC_Browse: not connected\n")
		return nil, errors.New("not connected")
	}

	ctx := context.Background()

	req := &mlibgrpc.BrowseRequest{
		Path:    path,
		Search:  search,
		Reverse: true,
	}

	resp, err := client.Browse(ctx, req)
	if err != nil {
		log.Printf("MLibGRPC_Browse: failed to browse: %v\n", err)
		return nil, err
	}

	log.Printf("MLibGRPC_Browse: %d results\n", len(resp.Items))

	return resp.Items, nil
}

func main() {

}
