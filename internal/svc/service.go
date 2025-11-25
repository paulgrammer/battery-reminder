package svc

import (
	"fmt"
	"log"

	"github.com/paulgrammer/battery-reminder/internal/config"
	"github.com/paulgrammer/battery-reminder/internal/monitor"

	"github.com/kardianos/service"
)

// ServiceManager handles system service operations
type ServiceManager struct {
	config  *config.Config
	service service.Service
	logger  service.Logger
}

// program implements the service.Interface
type program struct {
	config  *config.Config
	monitor *monitor.BatteryMonitor
	logger  service.Logger
}

// NewServiceManager creates a new service manager
func NewServiceManager(cfg *config.Config) (*ServiceManager, error) {
	svcConfig := &service.Config{
		Name:        "battery-reminder",
		DisplayName: "Battery Reminder Service",
		Description: "Monitors battery level and reminds when it's outside the optimal range",
		Arguments: []string{
			"-service", "run",
			"-min", fmt.Sprintf("%d", cfg.MinPercentage),
			"-max", fmt.Sprintf("%d", cfg.MaxPercentage),
		},
	}

	if cfg.SoundPath != "" {
		svcConfig.Arguments = append(svcConfig.Arguments, "-sound", cfg.SoundPath)
	}

	prg := &program{
		config: cfg,
	}

	svc, err := service.New(prg, svcConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	logger, err := svc.Logger(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	prg.logger = logger

	return &ServiceManager{
		config:  cfg,
		service: svc,
		logger:  logger,
	}, nil
}

// Control executes a service control action
func (sm *ServiceManager) Control(action string) error {
	if action == "run" {
		return sm.service.Run()
	}

	if err := service.Control(sm.service, action); err != nil {
		return fmt.Errorf("service control failed: %w", err)
	}

	return nil
}

// Start implements service.Interface
func (p *program) Start(s service.Service) error {
	p.logger.Info("Battery Reminder service starting.")
	p.logger.Info(p.config.String())

	p.monitor = monitor.NewBatteryMonitor(p.config)

	go p.monitor.Start()

	return nil
}

// Stop implements service.Interface
func (p *program) Stop(s service.Service) error {
	p.logger.Info("Battery Reminder service stopping.")
	return nil
}

// Install installs the service
func (sm *ServiceManager) Install() error {
	return sm.Control("install")
}

// Uninstall uninstalls the service
func (sm *ServiceManager) Uninstall() error {
	return sm.Control("uninstall")
}

// Start starts the service
func (sm *ServiceManager) Start() error {
	return sm.Control("start")
}

// Stop stops the service
func (sm *ServiceManager) Stop() error {
	return sm.Control("stop")
}

// Restart restarts the service
func (sm *ServiceManager) Restart() error {
	return sm.Control("restart")
}

// Run runs the service
func (sm *ServiceManager) Run() error {
	errs := make(chan error, 5)
	logger, err := sm.service.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	sm.logger = logger

	return sm.service.Run()
}
