package service

// DispatcherInterface defines the methods used by services to manage the
// lifecycle of service monitors.
type DispatcherInterface interface {
	AddService(serviceID uint)
	RemoveService(serviceID uint)
	RestartService(serviceID uint)
	Start()
	Stop()
}
