//go:build !windows

package gounrar

/*
#cgo CXXFLAGS: -g -ggdb -Wno-logical-op-parentheses -Wno-switch -Wno-dangling-else -D_FILE_OFFSET_BITS=64 -D_LARGEFILE_SOURCE -DRAR_SMP -DRARDLL

#include "dll.hpp"
*/
import "C"

func setOpenPath(option *C.RAROpenArchiveDataEx, path string) error {
	option.ArcName = C.CString(path)
	return nil
}

func getFileName(hdr *C.RARHeaderDataEx) string {
	return C.GoString(&hdr.FileName[0])
}
