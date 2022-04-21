package main

/*
#include <stdlib.h>
#include <stdint.h> // needed for windows
typedef const void constptr_t;
typedef const char cchar_t;
*/
import "C"
import (
	"log"

	"github.com/imperviousinc/beacon/components/core/internal"
)

//export BeaconHelper_Launch
//goland:noinspection GoSnakeCaseUsage
func BeaconHelper_Launch() C.int32_t {
	a, err := internal.NewAPI()
	if err != nil {
		log.Fatal(err)
	}
	a.Launch()
	return C.int32_t(0)
}

//export BeaconHelper_Shutdown
//goland:noinspection GoSnakeCaseUsage
func BeaconHelper_Shutdown() {
	log.Println("shutdown called this is cgo")
}

// TODO: add cert verification bindings here so it can be exposed via Mojo interface
// instead of using gRPC.

func main() {
	// empty main
}
