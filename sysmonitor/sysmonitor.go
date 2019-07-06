package sysmonitor

import (
	"unsafe"
)

var monitors []Monitor

type Monitor struct {
	hPhysicalMonitor uintptr
}

func (m *Monitor) Destroy() error {
	return destroyPhysicalMonitor(m.hPhysicalMonitor)
}

func (m *Monitor) SetBrightness(perc int) error {
	// TODO: convert from perc to value based on monitor
	return setMonitorBrightness(m.hPhysicalMonitor, uint32(perc))
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
		if supportsBrightness(monitor) {
			// TODO: implement get brightness call to get min and max values
			monitors = append(monitors, Monitor{
				hPhysicalMonitor: monitor.hPhysicalMonitor},
			)
		}
	}

	result := true
	return uintptr(unsafe.Pointer(&result))
}

func supportsBrightness(monitor physicalMonitor) bool {
	capabilities, supportedColorTemperatures := uint32(0), uint32(0)
	err := getMonitorCapabilities(monitor.hPhysicalMonitor, &capabilities, &supportedColorTemperatures)
	if err != nil {
		return false
	}

	if (capabilities & mcCapsBrightness) > 0 {
		return true
	}

	return false
}
