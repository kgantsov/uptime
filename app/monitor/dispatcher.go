package monitor

import (
	"sync"

	"github.com/kgantsov/uptime/app/repository"
	"github.com/sirupsen/logrus"
)

type Dispatcher struct {
	serviceRepo   repository.ServiceRepository
	heartbeatRepo repository.HeartbeatRepository
	monitors      map[uint]*Monitor
	mux           sync.Mutex
	logger        *logrus.Logger
}

func NewDispatcher(serviceRepo repository.ServiceRepository, heartbeatRepo repository.HeartbeatRepository, logger *logrus.Logger) *Dispatcher {
	d := &Dispatcher{
		serviceRepo:   serviceRepo,
		heartbeatRepo: heartbeatRepo,
		monitors:      make(map[uint]*Monitor),
		logger:        logger,
	}

	d.init()

	return d
}

func (d *Dispatcher) init() {
	services, err := d.serviceRepo.GetAll()
	if err != nil {
		d.logger.Errorf("Failed to load services: %s", err)
		return
	}

	d.mux.Lock()
	defer d.mux.Unlock()

	d.logger.Debugf("Found %d services", len(services))

	for i := range services {
		service := services[i]

		if service.Enabled {
			m := NewMonitor(d.heartbeatRepo, d.logger, &service)
			d.monitors[service.ID] = m
		}
	}
}

func (d *Dispatcher) AddService(serviceID uint) {
	d.logger.Infof("AddService %d", serviceID)

	service, err := d.serviceRepo.GetByID(serviceID)
	if err != nil {
		d.logger.Errorf("AddService: service %d not found: %s", serviceID, err)
		return
	}

	d.mux.Lock()
	defer d.mux.Unlock()

	if service.Enabled {
		m := NewMonitor(d.heartbeatRepo, d.logger, service)
		go m.Start()
		d.monitors[service.ID] = m
	}
}

func (d *Dispatcher) RemoveService(serviceID uint) {
	d.logger.Infof("RemoveService %d", serviceID)
	d.mux.Lock()
	defer d.mux.Unlock()

	monitor, ok := d.monitors[serviceID]
	if !ok {
		return
	}

	monitor.Stop()
	delete(d.monitors, serviceID)
}

func (d *Dispatcher) RestartService(serviceID uint) {
	d.RemoveService(serviceID)
	d.AddService(serviceID)
}

func (d *Dispatcher) Start() {
	for _, monitor := range d.monitors {
		go monitor.Start()
	}
}

func (d *Dispatcher) Stop() {
	for _, monitor := range d.monitors {
		monitor.Stop()
	}
}
