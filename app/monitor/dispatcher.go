package monitor

type Dispatcher struct {
	monitors []*Monitor
}

func NewDispatcher(services []Service) *Dispatcher {
	monitors := []*Monitor{}

	for _, service := range services {
		m := NewMonitor(service)
		monitors = append(monitors, m)
	}

	m := &Dispatcher{monitors: monitors}

	return m
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
