package common

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
)

func SetNotInput() {
	fd := uintptr(syscall.Stdin)

	// 获取控制台模式
	var mode uint32
	if err := windows.GetConsoleMode(windows.Handle(fd), &mode); err != nil {
		fmt.Println("Error getting console mode:", err)
		return
	}

	// 修改控制台模式，设置输入模式为不接受任何输入
	mode &^= windows.ENABLE_ECHO_INPUT | windows.ENABLE_LINE_INPUT | windows.ENABLE_PROCESSED_INPUT
	if err := windows.SetConsoleMode(windows.Handle(fd), mode); err != nil {
		fmt.Println("Error setting console mode:", err)
		return
	}
}
