package z9

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Module interface {
	Init() error
	Run()
	Stop()
}

type ModuleManager struct {
	modules []Module
}

func NewModuleManager() *ModuleManager {
	return &ModuleManager{
		modules: make([]Module, 0),
	}
}

func (m *ModuleManager) Append(mod Module) {
	m.modules = append(m.modules, mod)
}

func (m *ModuleManager) Init() error {
	for i := 0; i < len(m.modules); i++ {
		mod := m.modules[i]
		err := mod.Init()
		if err != nil {
			return fmt.Errorf("module manager init failed module=%v error=%v", mod, err)
		}
	}
	return nil
}

func (m *ModuleManager) Run() {
	for i := 0; i < len(m.modules); i++ {
		m.modules[i].Run()
	}
}

func (m *ModuleManager) Stop() {
	var wg sync.WaitGroup
	for i := 0; i < len(m.modules); i++ {
		wg.Add(1)
		go func(mod Module) {
			mod.Stop()
			wg.Done()
		}(m.modules[i])
	}
	wg.Wait()
}

func WaitExit() {
	exitChan := make(chan struct{})
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		close(exitChan)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan
}
