package conductrpc

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/viper"
)

// Entity represents client long polling instance.
type Entity struct {
	id     NodeID  // this
	client *client // RPC client

	state atomic.Value // the last state (can be nil first time)

	quitOnce sync.Once     //
	quit     chan struct{} //
}

// State returns current state.
func (e *Entity) State() (state *State) {
	if val := e.state.Load(); val != nil {
		return val.(*State)
	}
	return // nil, not polled yet
}

// SetState sets current state.
func (e *Entity) SetState(state *State) {
	e.state.Store(state) // update
}

// NewEntity creates RPC client for integration tests.
func NewEntity(id string) (e *Entity) {
	var (
		client, err = newClient(viper.GetString("integration_tests.address"))
		interval    = viper.GetDuration("integration_tests.lock_interval")
		state       *State
	)
	if err != nil {
		log.Fatalf("creating RPC client: %v", err)
	}

	e = new(Entity)
	e.id = NodeID(id)
	e.client = client
	e.quit = make(chan struct{})

	// initial state polling and wait node unlock
	for {
		// state polling can't return nil-State if err is nil
		state, err = client.state(NodeID(id))
		if err != nil {
			panic("requesting RPC (State): " + err.Error())
		}
		e.SetState(state)
		if !state.IsLock {
			break // can join blockchain
		}
		// otherwise, have to wait, retry after the interval

		// the joining is expected, since we can simply use the time.Sleep
		// instead of select with time.After and context.Done for tests
		time.Sleep(interval)
	}

	// start state polling
	go e.pollState()
	return
}

func (e *Entity) pollState() {
	for {
		select {
		case <-e.quit:
			return
		default:
		}
		var state, err = e.client.state(e.id)
		if err != nil {
			log.Printf("polling State: %v", err)
			continue
		}
		e.SetState(state)
	}
}

func (e *Entity) Shutdown() {
	e.quitOnce.Do(func() {
		close(e.quit)
	})
}

func (e *Entity) isMonitor() bool { //nolint:unused,deadcode // might be used later?
	var state = e.State()
	return state != nil && state.IsMonitor
}

//
// global
//

var global *Entity

// Init creates global Entity and locks until unlocked.
func Init(id string) {
	global = NewEntity(id)
}

// Shutdown the global Entity.
func Shutdown() {
	if global != nil {
		global.Shutdown()
	}
}

// Client returns global Entity to interact with. Use it, for example,
//
//     var state = conductrpc.Client().State()
//     for _, minerID := range miners {
//         if state.VRFS.IsBad(state, minerID) {
//             // send bad VRFS to this miner
//         } else if state.VRFS.IsGood(state, minerID) {
//             // send good VRFS to this miner
//         } else {
//             // don't send a VRFS to this miner
//         }
//     }
//
func Client() *Entity {
	return global
}
