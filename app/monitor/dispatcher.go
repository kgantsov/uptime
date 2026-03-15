package monitor

import (
	"sync"

	"github.com/kgantsov/uptime/app/repository"
	"github.com/rs/zerolog/log"
)

type Dispatcher struct {
	serviceRepo   repository.ServiceRepository
	heartbeatRepo repository.HeartbeatRepository
	monitors      map[uint]*Monitor
	mux           sync.Mutex
}

func NewDispatcher(serviceRepo repository.ServiceRepository, heartbeatRepo repository.HeartbeatRepository) *Dispatcher {
	d := &Dispatcher{
		serviceRepo:   serviceRepo,
		heartbeatRepo: heartbeatRepo,
		monitors:      make(map[uint]*Monitor),
	}

	d.init()

	return d
}

func (d *Dispatcher) init() {
	services, err := d.serviceRepo.GetAll()
	if err != nil {
		log.Error().Err(err).Msg("Failed to load services")
		return
	}

	d.mux.Lock()
	defer d.mux.Unlock()

	log.Debug().Int("count", len(services)).Msg("Found services")

	for i := range services {
		service := services[i]

		if service.Enabled {
			m := NewMonitor(d.heartbeatRepo, &service)
			d.monitors[service.ID] = m
		}
	}
}

func (d *Dispatcher) AddService(serviceID uint) {
	log.Info().Uint("serviceID", serviceID).Msg("AddService")

	service, err := d.serviceRepo.GetByID(serviceID)
	if err != nil {
		log.Error().Err(err).Uint("serviceID", serviceID).Msg("AddService: service not found")
		return
	}

	d.mux.Lock()
	defer d.mux.Unlock()

	if service.Enabled {
		m := NewMonitor(d.heartbeatRepo, service)
		go m.Start()
		d.monitors[service.ID] = m
	}
}

func (d *Dispatcher) RemoveService(serviceID uint) {
	log.Info().Uint("serviceID", serviceID).Msg("RemoveService")
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
