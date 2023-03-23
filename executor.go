package util

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

// HandlePanic logs goroutine panic by default
var HandlePanic = func(recovered interface{}, funcName string) {
	log.Println(fmt.Sprintf("%s panic: %v", funcName, recovered))
	log.Println(string(debug.Stack()))
}

// Executor is a executor without limits on counts of alive goroutines
// it tracks the goroutine started by it, and can cancel them when shutdown
type Executor struct {
	ctx                   context.Context
	cancel                context.CancelFunc
	activeGoroutinesMutex *sync.Mutex
	activeGoroutines      map[string]int
	HandlePanic           func(recovered interface{}, funcName string)
}

// GlobalExecutor has the life cycle of the program itself
// any goroutine want to be shutdown before main exit can be started from this executor
// GlobalExecutor expects the main function to call stop
// it does not magically knows the main function exits
var GlobalExecutor = NewExecutor()

// NewExecutor creates a new Executor,
// Executor can not be created by &Executor{}
// HandlePanic can be set with a callback to override global HandlePanic
func NewExecutor() *Executor {
	ctx, cancel := context.WithCancel(context.TODO())
	return &Executor{
		ctx:                   ctx,
		cancel:                cancel,
		activeGoroutinesMutex: &sync.Mutex{},
		activeGoroutines:      map[string]int{},
	}
}

// Go starts a new goroutine and tracks its lifecycle.
// Panic will be recovered and logged automatically, except for StopSignal
func (executor *Executor) Go(handler func(ctx context.Context)) {
	pc := reflect.ValueOf(handler).Pointer()
	f := runtime.FuncForPC(pc)
	funcName := f.Name()
	file, line := f.FileLine(pc)
	executor.activeGoroutinesMutex.Lock()
	defer executor.activeGoroutinesMutex.Unlock()
	startFrom := fmt.Sprintf("%s:%d", file, line)
	executor.activeGoroutines[startFrom] += 1
	go func() {
		defer func() {
			recovered := recover()
			// if you want to quit a goroutine without trigger HandlePanic
			// use runtime.Goexit() to quit
			if recovered != nil {
				if executor.HandlePanic == nil {
					HandlePanic(recovered, funcName)
				} else {
					executor.HandlePanic(recovered, funcName)
				}
			}
			executor.activeGoroutinesMutex.Lock()
			executor.activeGoroutines[startFrom] -= 1
			executor.activeGoroutinesMutex.Unlock()
		}()
		handler(executor.ctx)
	}()
}

// Stop cancel all goroutines started by this executor without wait
func (executor *Executor) Stop() {
	executor.cancel()
}

// StopAndWaitForever cancel all goroutines started by this executor and
// wait until all goroutines exited
func (executor *Executor) StopAndWaitForever() {
	executor.StopAndWait(context.Background())
}

// StopAndWait cancel all goroutines started by this executor and wait.
// Wait can be cancelled by the context passed in.
func (executor *Executor) StopAndWait(ctx context.Context) {
	executor.cancel()
	for {
		oneHundredMilliseconds := time.NewTimer(time.Millisecond * 100)
		select {
		case <-oneHundredMilliseconds.C:
			if executor.checkNoActiveGoroutines() {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (executor *Executor) checkNoActiveGoroutines() bool {
	executor.activeGoroutinesMutex.Lock()
	defer executor.activeGoroutinesMutex.Unlock()
	for startFrom, count := range executor.activeGoroutines {
		if count > 0 {
			log.Println("Executor is still waiting goroutines to quit",
				"startFrom", startFrom,
				"count", count)
			return false
		}
	}
	return true
}
