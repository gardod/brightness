// +build windows

package sysmonitor

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	u32                                      = windows.NewLazySystemDLL("user32.dll")
	d32                                      = windows.NewLazySystemDLL("dxva2.dll")
	pEnumDisplayMonitors                     = u32.NewProc("EnumDisplayMonitors")
	pGetNumberOfPhysicalMonitorsFromHMONITOR = d32.NewProc("GetNumberOfPhysicalMonitorsFromHMONITOR")
	pGetPhysicalMonitorsFromHMONITOR         = d32.NewProc("GetPhysicalMonitorsFromHMONITOR")
	pDestroyPhysicalMonitor                  = d32.NewProc("DestroyPhysicalMonitor")
	pGetMonitorCapabilities                  = d32.NewProc("GetMonitorCapabilities")
	pSetMonitorBrightness                    = d32.NewProc("SetMonitorBrightness")
)

const (
	mcCapsBrightness uint32 = 0x00000002
)

type physicalMonitor struct {
	hPhysicalMonitor             uintptr
	szPhysicalMonitorDescription [128]uint16
}

type rect struct {
	Left, Top, Right, Bottom int32
}

type monitorEnumDelegate func(hMonitor, hdc uintptr, lprcClip *rect, dwData uintptr) uintptr

func enumDisplayMonitors(hdc uintptr, lprcClip *rect, lpfnEnum monitorEnumDelegate, dwData uintptr) error {
	res, _, err := pEnumDisplayMonitors.Call(
		hdc,
		uintptr(unsafe.Pointer(lprcClip)),
		windows.NewCallback(lpfnEnum),
		dwData,
	)
	if res == 0 {
		return err
	}
	return nil
}

func getNumberOfPhysicalMonitorsFromHMONITOR(hMonitor uintptr, pdwNumberOfPhysicalMonitors *uint32) error {
	res, _, err := pGetNumberOfPhysicalMonitorsFromHMONITOR.Call(
		hMonitor,
		uintptr(unsafe.Pointer(pdwNumberOfPhysicalMonitors)),
	)
	if res == 0 {
		return err
	}
	return nil
}

func getPhysicalMonitorsFromHMONITOR(hMonitor uintptr, dwPhysicalMonitorArraySize uint32, pPhysicalMonitorArray []physicalMonitor) error {
	res, _, err := pGetPhysicalMonitorsFromHMONITOR.Call(
		hMonitor,
		uintptr(dwPhysicalMonitorArraySize),
		uintptr(unsafe.Pointer(&pPhysicalMonitorArray[0])),
	)
	if res == 0 {
		return err
	}
	return nil
}

func destroyPhysicalMonitor(hMonitor uintptr) error {
	res, _, err := pDestroyPhysicalMonitor.Call(
		hMonitor,
	)
	if res == 0 {
		return err
	}
	return nil
}

func getMonitorCapabilities(hMonitor uintptr, pdwMonitorCapabilities, pdwSupportedColorTemperatures *uint32) error {
	res, _, err := pGetMonitorCapabilities.Call(
		hMonitor,
		uintptr(unsafe.Pointer(pdwMonitorCapabilities)),
		uintptr(unsafe.Pointer(pdwSupportedColorTemperatures)),
	)
	if res == 0 {
		return err
	}
	return nil
}

func setMonitorBrightness(hMonitor uintptr, dwNewBrightness uint32) error {
	res, _, err := pSetMonitorBrightness.Call(
		hMonitor,
		uintptr(dwNewBrightness),
	)
	if res == 0 {
		return err
	}
	return nil
}
