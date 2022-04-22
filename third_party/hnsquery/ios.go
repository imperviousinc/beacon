//go:build ios
// +build ios

package hnsquery

/*
#cgo CFLAGS: -xobjective-c -fmodules -fobjc-arc
#import <Foundation/Foundation.h>
static inline void activate_cocoa_multithreading() {
    [[NSThread new] start];
}
#pragma GCC visibility push(hidden)
*/
import "C"

func init() {
	C.activate_cocoa_multithreading()
}
