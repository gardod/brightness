package sysmonitor

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var monitors []Monitor

type Monitor struct {
	hPhysicalMonitor uintptr

	Name          string
	MinBrightness int64
	MaxBrightness int64
}

func (m *Monitor) Destroy() error {
	return destroyPhysicalMonitor(m.hPhysicalMonitor)
}

func (m *Monitor) SetBrightness(perc int64) error {
	value := m.MinBrightness + (m.MaxBrightness-m.MinBrightness)*perc/100
	return setMonitorBrightness(m.hPhysicalMonitor, uint32(value))
}

func GetMonitors() ([]Monitor, error) {
	err := enumDisplayMonitors(0, nil, monitorEnumDelegate(getPhysicalMonitors), 0)
	if err != nil {
		return nil, err
	}

	result := monitors
	monitors = nil
	return result, nil
}

func getPhysicalMonitors(hMonitor, hdc uintptr, lprcClip *rect, dwData uintptr) uintptr {
	monitorCount := uint32(0)
	err := getNumberOfPhysicalMonitorsFromHMONITOR(hMonitor, &monitorCount)
	if err != nil {
		result := false
		return uintptr(unsafe.Pointer(&result))
	}
	if monitorCount <= 0 {
		result := true
		return uintptr(unsafe.Pointer(&result))
	}

	physicalMonitors := make([]physicalMonitor, monitorCount)
	err = getPhysicalMonitorsFromHMONITOR(hMonitor, monitorCount, physicalMonitors)
	if err != nil {
		result := false
		return uintptr(unsafe.Pointer(&result))
	}

	for _, monitor := range physicalMonitors {
		addMonitor(monitor)
	}

	result := true
	return uintptr(unsafe.Pointer(&result))
}

func addMonitor(monitor physicalMonitor) {
	if !supportsBrightness(monitor) {
		return
	}

	minBrightness, curBrightness, maxBrightness := uint32(0), uint32(0), uint32(0)
	err := getMonitorBrightness(monitor.hPhysicalMonitor, &minBrightness, &curBrightness, &maxBrightness)
	if err != nil {
		return
	}

	monitors = append(monitors, Monitor{
		hPhysicalMonitor: monitor.hPhysicalMonitor,
		Name:             windows.UTF16ToString(monitor.szPhysicalMonitorDescription[:]),
		MinBrightness:    int64(minBrightness),
		MaxBrightness:    int64(maxBrightness),
	})
}

func supportsBrightness(monitor physicalMonitor) bool {
	capabilities, supportedColorTemperatures := uint32(0), uint32(0)
	err := getMonitorCapabilities(monitor.hPhysicalMonitor, &capabilities, &supportedColorTemperatures)
	if err != nil {
		return false
	}

	return (capabilities & mcCapsBrightness) > 0
}
