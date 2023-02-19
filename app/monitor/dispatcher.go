package monitor

import (
	"sync"

	"github.com/kgantsov/uptime/app/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Dispatcher struct {
	DB       *gorm.DB
	monitors map[uint]*Monitor
	mux      sync.Mutex
	logger   *logrus.Logger
}

func NewDispatcher(db *gorm.DB, logger *logrus.Logger) *Dispatcher {
	d := &Dispatcher{
		DB:       db,
		monitors: make(map[uint]*Monitor),
		logger:   logger,
	}

	d.init()

	return d
}

func (d *Dispatcher) init() {
	services := d.getServices()

	d.mux.Lock()
	defer d.mux.Unlock()

	d.logger.Debugf("Found %d services", len(services))

	for i := range services {
		service := services[i]

		if service.Enabled {
			m := NewMonitor(d.DB, d.logger, &service)
			d.monitors[service.ID] = m
		}
	}
}

func (d *Dispatcher) getServices() []model.Service {
	var services []model.Service

	err := d.DB.Model(&model.Service{}).Preload("Notifications").Order("id desc").Find(&services).Error

	if err != nil {
		return services
	}

	return services
}

func (d *Dispatcher) AddService(serviceID uint) {
	d.logger.Infof("AddService %d", serviceID)

	var service model.Service

	err := d.DB.Model(&model.Service{}).Preload("Notifications").Order("id desc").First(&service).Error

	if err != nil {
		return
	}

	d.mux.Lock()
	defer d.mux.Unlock()

	if service.Enabled {
		m := NewMonitor(d.DB, d.logger, &service)
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
