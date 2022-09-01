package gounrar

/*
#cgo CXXFLAGS: -g -ggdb -Wno-logical-op-parentheses -Wno-switch -Wno-dangling-else -D_FILE_OFFSET_BITS=64 -D_LARGEFILE_SOURCE -DRAR_SMP -DRARDLL

#include <stdlib.h>
#include "dll.hpp"

int callbackInC(UINT msg,uintptr_t UserData,uintptr_t P1,uintptr_t P2);
*/
import "C"
import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"runtime/cgo"
	"time"
	"unsafe"
)

var (
	ErrIsDirectory = errors.New("is directory")
)

type Error struct {
	code    int
	message string
}

func (o Error) Error() string {
	return fmt.Sprintf("%s: rar error(%d)", o.message, o.code)
}

type Header struct {
	FileName   string
	Type       uint
	Flags      uint // directory or file
	PackSize   int64
	UnpackSize int64
	CRC        uint
	ModTime    time.Time
	CreateTime time.Time
	AccessTime time.Time
	BlockPos   int64
}

func (h *Header) IsDir() bool {
	return (h.Flags & uint(C.RHDF_DIRECTORY)) != 0
}

func (h *Header) IsSolid() bool {
	return (h.Flags & uint(C.RHDF_SOLID)) != 0
}

func (h *Header) IsEncrypted() bool {
	return (h.Flags & uint(C.RHDF_ENCRYPTED)) != 0
}

type Archive struct {
	path              string
	ptr               unsafe.Pointer
	currentHeader     *Header // repeat read cause BAD_DATA, don't know why
	currentFileBuffer *bytes.Buffer
}

func (a *Archive) Close() error {
	a.currentHeader = nil
	a.currentFileBuffer = nil
	code := C.RARCloseArchive(a.ptr)
	a.ptr = nil
	if code == C.ERAR_SUCCESS {
		return nil
	} else {
		return Error{message: "close", code: int(code)}
	}
}

func (a *Archive) GetHeader() (*Header, error) {
	if a.currentHeader != nil {
		return a.currentHeader, nil
	}
	return a.readHeader()
}

func (a *Archive) readHeader() (*Header, error) {
	hdr := C.RARHeaderDataEx{}
	code := C.RARReadHeaderEx(a.ptr, &hdr)
	if code != C.ERAR_SUCCESS {
		if code == C.ERAR_END_ARCHIVE {
			return nil, io.EOF
		}
		return nil, Error{message: "readhead", code: int(code)}
	}

	header := Header{
		FileName:   getFileName(&hdr),
		Type:       uint(hdr.Type),
		Flags:      uint(hdr.Flags),
		PackSize:   int64(hdr.PackSize),
		UnpackSize: makeInt64(hdr.UnpSizeHigh, hdr.UnpSize),
		CRC:        uint(hdr.FileCRC),
		ModTime:    convertTime(hdr.MtimeUnix),
		CreateTime: convertTime(hdr.CtimeUnix),
		AccessTime: convertTime(hdr.AtimeUnix),
		BlockPos:   int64(hdr.BlockPos),
	}

	a.currentHeader = &header
	return &header, nil
}

// Next move to next entry and read head, so it can tell us whether EOF
func (a *Archive) Next() (*Header, error) {
	if a.currentHeader != nil {
		// skip current file
		_, err := a.processFile(C.RAR_SKIP, 0)
		if err != nil {
			return nil, err
		}
	}

	return a.readHeader()
}

func (a *Archive) SeekPos(blockPos int64) (*Header, error) {
	if err := a.seek(blockPos); err != nil {
		return nil, err
	}

	a.currentFileBuffer = nil
	return a.currentHeader, nil
}

func (a *Archive) seek(blockPos int64) error {
	code := C.RARSeek(a.ptr, C.int64_t(blockPos))
	if code != C.ERAR_SUCCESS {
		return Error{message: "seek", code: int(code)}
	}

	_, err := a.readHeader() // readHeader after seek so ReadAll can work
	return err
}

// ReadAll read all content of current file
func (a *Archive) ReadAll() ([]byte, error) {
	hdr, err := a.GetHeader()
	if err != nil {
		return nil, err
	}

	if hdr.IsDir() {
		return nil, ErrIsDirectory
	}

	// read all file content to buffer
	buf, err := a.processFile(C.RAR_TEST, hdr.UnpackSize)
	if err != nil {
		return nil, err
	}

	a.currentFileBuffer = buf
	// rewind file cursor so that ReadAll again return same content
	if err := a.seek(hdr.BlockPos); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (a *Archive) processFile(op C.int, initalBufferSize int64) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(make([]byte, 0, initalBufferSize))
	h := cgo.NewHandle(buf)
	defer h.Delete()
	C.RARSetCallback(a.ptr, (C.UNRARCALLBACK)(unsafe.Pointer(C.callbackInC)), C.uintptr_t(h))
	code := C.RARProcessFile(a.ptr, op, nil, nil) // USE RAR_TEST to avoid write file
	C.RARSetCallback(a.ptr, nil, C.uintptr_t(0))  // reset callback
	if code != C.ERAR_SUCCESS {
		return nil, Error{message: "process", code: int(code)}
	}

	return buf, nil
}

// Read current file in Archive
func (a *Archive) Read(p []byte) (int, error) {
	if a.currentFileBuffer == nil {
		_, err := a.ReadAll()
		if err != nil {
			return 0, err
		}
	}

	return a.currentFileBuffer.Read(p)
}

type Options struct {
}

type Option func(options *Options)

func Open(path string, options ...Option) (*Archive, error) {
	var opt Options
	for _, apply := range options {
		apply(&opt)
	}
	option := C.RAROpenArchiveDataEx{}
	if err := setOpenPath(&option, path); err != nil {
		return nil, err
	}
	if option.ArcName != nil {
		defer C.free(unsafe.Pointer(option.ArcName))
	}

	if option.ArcNameW != nil {
		defer C.free(unsafe.Pointer(option.ArcNameW))
	}

	option.OpenMode = C.RAR_OM_EXTRACT // or RAR_OM_LIST
	hdl := C.RAROpenArchiveEx(&option)
	if C.ERAR_SUCCESS == option.OpenResult {
		return &Archive{
			path: path,
			ptr:  hdl,
		}, nil
	} else {
		return nil, Error{message: "open", code: int(option.OpenResult)}
	}
}

func makeInt64(high, low C.uint) int64 {
	val := int64(high)<<32 | int64(low)
	return val
}

func convertTime(t C.uint64_t) time.Time {
	ns := uint64(t)
	return time.Unix(int64(ns/1e9), int64(ns%1e9))
}

// callbackInGo return -1 cause error 21
//export callbackInGo
func callbackInGo(msg int, user, p1, p2 C.uintptr_t) C.int {
	if msg != int(C.UCM_PROCESSDATA) {
		return -1
	}
	a := cgo.Handle(user).Value().(*bytes.Buffer)
	//goland:noinspection GoVetUnsafePointer
	ptr := unsafe.Pointer(uintptr(p1))
	data := unsafe.Slice((*byte)(ptr), uint64(p2))
	n, err := a.Write(data)
	if err != nil {
		panic(fmt.Errorf("write to buf error: %s", err))
		return -1
	}
	if n != len(data) {
		panic(fmt.Errorf("write wrong len(data)=%d, write=%d", len(data), n))
		return -1
	}
	return 0
}
