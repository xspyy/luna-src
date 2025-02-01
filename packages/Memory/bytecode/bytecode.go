package bytecode

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"syscall"
	"unsafe"

	"github.com/cespare/xxhash/v2"
	"github.com/klauspost/compress/zstd"
	lua "github.com/yuin/gopher-lua"
)

var (
	dll               = syscall.NewLazyDLL("Luna.dll")
	procRBXCompile    = dll.NewProc("getbytecode")
	procRBXDecompress = dll.NewProc("decompress")
	loadErr           error
)

func init() {
	loadErr = dll.Load()
}

type Bytecode struct{}

var limiter sync.Mutex

func MBs(i int) int {
	return i * 1024 * 1024
}

/*
var Bytes []byte = make([]byte, MBs(10))

	func cloneV2(b []byte) []byte {
		if b == nil {
			return nil
		}
		nb := make([]byte, len(b))
		copy(nb, b)
		return nb
	}
*/
func (b *Bytecode) Compile(source string) ([]byte, int64) {

	if loadErr != nil {
		return []byte("Luna.DLL Error: " + loadErr.Error()), -10
	}

	limiter.Lock()
	defer limiter.Unlock()

	buffer := make([]byte, MBs(10))

	var actualSize uintptr
	procRBXCompile.Call(
		uintptr(unsafe.Pointer(&[]byte(source)[0])),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&actualSize)),
	)

	trimmedSize := len(buffer)
	for trimmedSize > 0 && buffer[trimmedSize-1] == 0 {
		trimmedSize--
	}

	return buffer[:trimmedSize], int64(trimmedSize)

}

func CompileTest() []byte {
	L := lua.NewState()
	defer L.Close()

	code := `return function()
    local x = 10
    return x * 2
end`

	fn, _ := L.LoadString(code)

	var bytecodeBuffer bytes.Buffer
	for _, instr := range fn.Proto.Code {
		binary.Write(&bytecodeBuffer, binary.LittleEndian, instr)
	}

	bytecodeBytes := bytecodeBuffer.Bytes()
	return bytecodeBytes
}

func Compress(bytecode []byte) ([]byte, int, error) {
	dataSize := len(bytecode)

	buffer := make([]byte, 8)

	copy(buffer[0:], "RSB1")

	binary.LittleEndian.PutUint32(buffer[4:], uint32(dataSize))

	encoder, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create ZSTD encoder: %v", err)
	}

	compressedData := encoder.EncodeAll(bytecode, nil)
	compressedSize := len(compressedData)

	if compressedSize == 0 {
		return nil, 0, errors.New("failed to compress the bytecode")
	}

	buffer = append(buffer, compressedData...)

	size := compressedSize + 8

	hasher := xxhash.New()
	hasher.Write(buffer[:size])
	hashBytes := hasher.Sum(nil)
	key := binary.LittleEndian.Uint32(hashBytes)

	keyBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(keyBytes, key)

	for i := 0; i < size; i++ {
		buffer[i] ^= keyBytes[i%4] + byte(i*41)
	}

	return buffer[:size], size, nil
}

func (b *Bytecode) Decompress(source []byte) []byte {

	if loadErr != nil {
		return []byte("Luna.DLL Error: " + loadErr.Error())
	}

	initialCapacity := 1024
	buffer := make([]byte, 0, initialCapacity)
	buffer = buffer[:initialCapacity]
	procRBXDecompress.Call(
		uintptr(unsafe.Pointer(&[]byte(source)[0])),
		uintptr(unsafe.Pointer(&buffer[0])),
	)
	return buffer
}
