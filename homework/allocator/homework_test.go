package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

// Defragment - Дефрагментирует memory, имзеняет указатели.
// Предполагается, что:
//  1. pointers ссылаются только на 1 байт памяти.
//  2. pointers отсортирован по возрастанию адреса.
//  3. len(pointers) <= len(memory).
func Defragment(memory []byte, pointers []unsafe.Pointer) {
	if len(memory) == 0 || len(pointers) == 0 {
		return
	}

	for i := range pointers {
		nextByte := unsafe.Pointer(&memory[i])
		swapPointersAndValues(&nextByte, &pointers[i])
	}
}

func swapPointersAndValues(a, b *unsafe.Pointer) {
	abyte := (*byte)(*a)
	bbyte := (*byte)(*b)

	*abyte, *bbyte = *bbyte, *abyte
	*a, *b = *b, *a
}

func TestDefragmentation(t *testing.T) {
	var fragmentedMemory = []byte{
		0xFF, 0x00, 0x00, 0x00,
		0x00, 0xFF, 0x00, 0x00,
		0x00, 0x00, 0xFF, 0x00,
		0x00, 0x00, 0x00, 0xFF,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[15]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
	}

	var defragmentedMemory = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}
