package main

/*
typedef struct MLibGRPC_BrowseItem {
	char *name;
	char *uri;
	char *image_uri;
	int folder;
} MLibGRPC_BrowseItem;

typedef struct MLibGRPC_BrowseItems {
	MLibGRPC_BrowseItem **items;
	int count;
} MLibGRPC_BrowseItems;

typedef struct MLibGRPC_MediaItems {
	char **items;
	int count;
} MLibGRPC_MediaItems;


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
func MLibGRPC_Browse(cURI *C.char, cSearch *C.char, browseType int32) *C.struct_MLibGRPC_BrowseItems {
	uri := C.GoString(cURI)
	search := C.GoString(cSearch)
	browseTypeName := mlibgrpc.BrowseType_name[browseType]

	log.Printf("MLibGRPC_Browse: uri: %s search: %s type: %s\n", uri, search, browseTypeName)

	items, err := browse(uri, search, mlibgrpc.BrowseType(browseType))
	if err != nil {
		log.Printf("MLibGRPC_Browse: error %v\n", err)
		items = []*mlibgrpc.BrowseItem{
			{
				Name: err.Error(),
			},
		}
	}

	result := C.malloc(C.size_t(len(items)) * C.size_t(unsafe.Sizeof(uintptr(0))))
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

	log.Printf("MLibGRPC_Browse: %d results\n", len(items))

	cBrowseItems := (*C.struct_MLibGRPC_BrowseItems)(C.malloc(
		C.size_t(unsafe.Sizeof(C.struct_MLibGRPC_BrowseItems{}))))
	cBrowseItems.items = (**C.struct_MLibGRPC_BrowseItem)(result)
	cBrowseItems.count = C.int(len(items))

	return (*C.struct_MLibGRPC_BrowseItems)(cBrowseItems)
}

func browse(uri string, search string, browseType mlibgrpc.BrowseType) ([]*mlibgrpc.BrowseItem, error) {
	connMutex.Lock()
	defer connMutex.Unlock()

	if client == nil {
		return nil, errors.New("not connected")
	}

	ctx := context.Background()

	req := &mlibgrpc.BrowseRequest{
		Uri:        uri,
		Search:     search,
		BrowseType: browseType,
	}

	resp, err := client.Browse(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

//export MLibGRPC_Media
func MLibGRPC_Media(cURI *C.char, cSearch *C.char, browseType int32) *C.struct_MLibGRPC_MediaItems {
	uri := C.GoString(cURI)
	search := C.GoString(cSearch)
	browseTypeName := mlibgrpc.BrowseType_name[browseType]

	log.Printf("MLibGRPC_Media: uri: %s type: %s\n", uri, browseTypeName)

	items, err := media(uri, search, mlibgrpc.BrowseType(browseType))
	if err != nil {
		log.Printf("MLibGRPC_Media: error %v\n", err)
		items = []string{err.Error()}
	}

	result := C.malloc(C.size_t(len(items)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	resultArr := (*[1<<30 - 1]*C.char)(result)

	for i, item := range items {
		resultArr[i] = C.CString(item)
	}

	log.Printf("MLibGRPC_Media: %d results\n", len(items))

	cMediaItems := (*C.struct_MLibGRPC_MediaItems)(C.malloc(
		C.size_t(unsafe.Sizeof(C.struct_MLibGRPC_MediaItems{}))))
	cMediaItems.items = (**C.char)(result)
	cMediaItems.count = C.int(len(items))

	return (*C.struct_MLibGRPC_MediaItems)(cMediaItems)
}

func media(uri string, search string, browseType mlibgrpc.BrowseType) ([]string, error) {
	connMutex.Lock()
	defer connMutex.Unlock()

	if client == nil {
		return nil, errors.New("not connected")
	}

	ctx := context.Background()

	req := &mlibgrpc.MediaRequest{
		Uri:        uri,
		Search:     search,
		BrowseType: browseType,
	}

	resp, err := client.Media(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Uris, nil
}

func main() {

}
