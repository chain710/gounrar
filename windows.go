//go:build windows

package gounrar

/*
#cgo LDFLAGS: -lwbemuuid -lole32 -loleaut32
#cgo CXXFLAGS: -g -ggdb -Wno-logical-op-parentheses -Wno-switch -Wno-dangling-else -D_FILE_OFFSET_BITS=64 -D_LARGEFILE_SOURCE -DRAR_SMP -DRARDLL

#include "dll.hpp"
*/
import "C"
import (
	"golang.org/x/sys/windows"
)

func WCharPtrToString(p *C.wchar_t) string {
	return windows.UTF16PtrToString((*uint16)(p))
}

func WCharPtrFromString(s string) (*C.wchar_t, error) {
	p, err := windows.UTF16PtrFromString(s)
	return (*C.wchar_t)(p), err
}

// return defer func and error
func setOpenPath(option *C.RAROpenArchiveDataEx, path string) error {
	wstr, err := WCharPtrFromString(path)
	if err != nil {
		return err
	}

	option.ArcNameW = C.wcsdup(wstr) // to avoid cgo pointer issue
	option.ArcName = nil
	return nil
}

func getFileName(hdr *C.RARHeaderDataEx) string {
	return WCharPtrToString(&hdr.FileNameW[0])
}
