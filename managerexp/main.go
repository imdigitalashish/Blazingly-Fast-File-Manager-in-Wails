package main

import (
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

func getSizeDrive(driveLetter string) (uint64, error) {
	var freeBytes Bytex
	var totalBytes Bytex
	var totalFreeBytes Bytex

	drivePath := syscall.StringToUTF16Ptr(driveLetter + "\\")
	ret, _, err := getDiskFreeSpaceExPtr.Call(uintptr(unsafe.Pointer(drivePath)),
		uintptr(unsafe.Pointer(&freeBytes)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&totalFreeBytes)))

	if ret == 0 {
		return 0, fmt.Errorf("Failed %w", err)
	}

	return uint64(totalBytes.QuadPart), nil
}

func getDrives() []string {
	var drives []string

	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXZ" {
		_, err := os.Open(string(drive) + ":\\")
		if err == nil {
			drives = append(drives, string(drive))
		}
	}

	driveSize, err := getSizeDrive("C:")
	if err != nil {
		fmt.Println("Error getting drive size: ", err)

	}

	fmt.Print("Size %d MB\n", driveSize/1024/1024/1024)

	return drives
}

func main() {

	fmt.Println(getDrives())
}
