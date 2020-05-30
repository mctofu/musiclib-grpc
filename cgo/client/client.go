package main

/*
typedef struct MLibGRPC_BrowseItem {
	char *name;
	char *uri;
	char *image_uri;
	int folder;
} MLibGRPC_BrowseItem;
*/
import "C"

import (
	"context"
	"errors"
	"log"
	"sync"
	"unsafe"

	"github.com/mctofu/music-library-grpc/go/mlibgrpc"
	"google.golang.org/grpc"
)

var client mlibgrpc.MusicLibraryClient
var conn *grpc.ClientConn
var connMutex sync.Mutex

//export MLibGRPC_Connect
func MLibGRPC_Connect() int {
	connMutex.Lock()
	defer connMutex.Unlock()

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
	connMutex.Lock()
	defer connMutex.Unlock()

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
func MLibGRPC_Browse(cPath *C.char, cSearch *C.char, browseType int32) **C.struct_MLibGRPC_BrowseItem {
	path := C.GoString(cPath)
	search := C.GoString(cSearch)

	browseTypeName := mlibgrpc.BrowseType_name[browseType]

	log.Printf("MLibGRPC_Browse: Browse path: %s search: %s type: %s\n", path, search, browseTypeName)

	items, err := browse(path, search, mlibgrpc.BrowseType(browseType))
	if err != nil {
		log.Printf("MLibGRPC_Browse: error %v\n", err)
		items = []*mlibgrpc.BrowseItem{
			{
				Name: err.Error(),
			},
		}
	}

	result := C.malloc(C.size_t(len(items)+1) * C.size_t(unsafe.Sizeof(uintptr(0))))
	resultArr := (*[1<<30 - 1]*C.struct_MLibGRPC_BrowseItem)(result)

	for i, item := range items {
		cBrowseItem := (*C.struct_MLibGRPC_BrowseItem)(C.malloc(
			C.size_t(unsafe.Sizeof(C.struct_MLibGRPC_BrowseItem{}))))
		cBrowseItem.name = C.CString(item.Name)
		cBrowseItem.uri = C.CString(item.Uri)
		cBrowseItem.image_uri = C.CString(item.ImageUri)
		if item.Folder {
			cBrowseItem.folder = 1
		} else {
			cBrowseItem.folder = 0
		}

		resultArr[i] = cBrowseItem
	}

	resultArr[len(items)] = nil

	log.Printf("MLibGRPC_Browse: %d results\n", len(items))

	return (**C.struct_MLibGRPC_BrowseItem)(result)
}

func browse(path string, search string, browseType mlibgrpc.BrowseType) ([]*mlibgrpc.BrowseItem, error) {
	connMutex.Lock()
	defer connMutex.Unlock()

	if client == nil {
		return nil, errors.New("not connected")
	}

	ctx := context.Background()

	req := &mlibgrpc.BrowseRequest{
		Path:       path,
		Search:     search,
		Reverse:    true,
		BrowseType: browseType,
	}

	resp, err := client.Browse(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

func main() {

}
