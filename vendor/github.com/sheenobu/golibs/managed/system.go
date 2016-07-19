package managed

import (
	"github.com/sheenobu/golibs/log"
	"golang.org/x/net/context"

	"sync"
	"time"
)

// A System is a group of related goroutines and subsystems
type System struct {
	name string

	parentContext context.Context
	parentCancel  func()

	ctx      context.Context
	cancelFn func()
	wg       sync.WaitGroup

	startChan chan string
	stopChan  chan string

	lock          sync.RWMutex
	Children      []*System
	ChildrenProcs []*Process
}

// NewSystem creates a new system
func NewSystem(name string) *System {
	return &System{
		name:      name,
		startChan: make(chan string),
		stopChan:  make(chan string),
	}
}

// Context returns the context of the application
func (sys *System) Context() context.Context {
	return sys.ctx
}

// Start starts the application
func (sys *System) Start() {

	sys.parentContext = context.Background()
	sys.parentContext, sys.parentCancel = context.WithCancel(sys.parentContext)

	sys.ctx, sys.cancelFn = context.WithCancel(sys.parentContext)

	log.Log(sys.ctx).Debug("Starting system", "name", sys.name)
}

// StartWithContext starts the application
func (sys *System) StartWithContext(ctx context.Context) {

	sys.parentContext = ctx
	sys.parentContext, sys.parentCancel = context.WithCancel(sys.parentContext)

	sys.ctx, sys.cancelFn = context.WithCancel(sys.parentContext)

	log.Log(sys.ctx).Debug("Starting system", "system", sys.name)
}

// StartWithParent starts the application with the parent sys as the context
func (sys *System) StartWithParent(parent *System) {

	sys.parentContext = context.Background()
	sys.parentContext, sys.parentCancel = context.WithCancel(sys.parentContext)

	sys.ctx, sys.cancelFn = context.WithCancel(parent.ctx)

	log.Log(sys.ctx).Debug("Starting system with parent system", "name", sys.name, "parent", parent.name)
}

// Stop stops the application
func (sys *System) Stop() {

	log.Log(sys.ctx).Debug("Stopping system", "name", sys.name)

	sys.cancelFn()
}

// Add the process to the system
func (sys *System) Add(process *Process) {
	go func() {
		sys.startChan <- process.Name
		sys.lock.Lock()
		sys.ChildrenProcs = append(sys.ChildrenProcs, process)
		sys.lock.Unlock()
		process.Set("Status", "Started")
		process.Component(sys.ctx)
		process.Set("Status", "Stopped")
		sys.stopChan <- process.Name
	}()
}

// SpawnSystem starts the system as a child
func (sys *System) SpawnSystem(child *System) {
	go func() {
		sys.startChan <- "app:" + child.name
		sys.lock.Lock()
		sys.Children = append(sys.Children, child)
		sys.lock.Unlock()
		child.StartWithParent(sys)
		child.Wait()
		sys.stopChan <- "app:" + child.name
	}()
}

// Wait waits for the application and its subprocesses to stop
func (sys *System) Wait() {
	log.Log(sys.ctx).Debug("Waiting on application stop", "app", sys.name)

	procs := make(map[string]bool)

	sys.watchProcessState(func(name string, status bool) {
		procs[name] = status
	})

	ch := make(chan struct{})

	log.Log(sys.ctx).Debug("System stopped, waiting on subprocesses", "app", sys.name)
	go func() {
		sys.wg.Wait()
		close(ch)
	}()

	select {
	case <-ch:
		log.Log(sys.ctx).Debug("All subprocesses stopped", "app", sys.name)
	case <-time.After(1 * time.Second):
		log.Log(sys.ctx).Error("Some subprocesses failed to stop in time", "app", sys.name)
		for k, v := range procs {
			if v {
				log.Log(sys.ctx).Error("Subprocess failed to stop in time", "proc", k, "app", sys.name)
			}
		}
	}

	sys.parentCancel()
}

func (sys *System) watchProcessState(setProcStatus func(name string, status bool)) {
	go func() {

	L:
		for {
			select {
			case proc := <-sys.startChan:
				log.Log(sys.ctx).Debug("Got subrocess start", "process", proc, "app", sys.name)
				sys.wg.Add(1)
				setProcStatus(proc, true)
			case proc := <-sys.stopChan:
				log.Log(sys.ctx).Debug("Got subprocess stop", "process", proc, "app", sys.name)
				sys.wg.Done()
				setProcStatus(proc, false)
			case <-sys.ctx.Done():
				break L
			}
		}

		for {
			select {
			case proc := <-sys.stopChan:
				log.Log(sys.ctx).Debug("Got subprocess stop", "process", proc, "app", sys.name)
				setProcStatus(proc, false)
				sys.wg.Done()
			case <-sys.parentContext.Done():
				return
			}
		}
	}()

	<-sys.ctx.Done()
}
