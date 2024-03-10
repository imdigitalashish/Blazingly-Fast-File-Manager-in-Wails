package main

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type Bytex struct {
	QuadPart int64
}

var (
	kernel32dll           = syscall.NewLazyDLL("kernel32.dll")
	getDiskFreeSpaceExPtr = kernel32dll.NewProc("GetDiskFreeSpaceExW")
)

func getSizeDrive(driveLetter string) ([]uint64, error) {
	var freeBytes Bytex
	var totalBytes Bytex
	var totalFreeBytes Bytex

	drivePath := syscall.StringToUTF16Ptr(driveLetter + "\\")
	ret, _, err := getDiskFreeSpaceExPtr.Call(uintptr(unsafe.Pointer(drivePath)),
		uintptr(unsafe.Pointer(&freeBytes)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&totalFreeBytes)))

	if ret == 0 {
		return nil, fmt.Errorf("Failed %w", err)
	}

	return []uint64{uint64(totalBytes.QuadPart), uint64(totalFreeBytes.QuadPart)}, nil
}

func getDrives() string {

	var drives = make(map[string]map[string]uint64)

	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXZ" {
		_, err := os.Open(string(drive) + ":\\")
		if err == nil {
			driveSize, err := getSizeDrive(string(drive) + ":")
			if err == nil {
				drives[string(drive)] = map[string]uint64{
					"total_space": driveSize[0] / 1024 / 1024 / 1024,
					"space_left":  driveSize[1] / 1024 / 1024 / 1024, // Assuming we can't get free space here (modify if possible)
				}
			}
		}
	}

	jsonData, err := json.Marshal(drives)
	if err != nil {
		fmt.Println("Error Marshalling map to JSON: ", err)
	}

	return string(jsonData)
}
