package ib

import (
	"errors"
	"flag"
	"os"
	"reflect"
	"testing"
	"time"
)

var (
	gwURL         = flag.String("gw", "", "Gateway URL")
	noEngineReuse = flag.Bool("no-engine-reuse", false, "Don't keep reusing the engine; each test case gets its own engine.")
)

func getGatewayURL() string {
	if *gwURL != "" {
		return *gwURL
	}
	if url := os.Getenv("GATEWAY_URL"); url != "" {
		return url
	}
	return "localhost:4002"
}

func (e *Engine) expect(t *testing.T, seconds int, ch chan Reply, expected []IncomingMessageID) (Reply, error) {
	for {
		select {
		case <-time.After(time.Duration(seconds) * time.Second):
			return nil, errors.New("Timeout waiting")
		case v := <-ch:
			if v.code() == 0 {
				t.Fatalf("don't know message '%v'", v)
			}
			for _, code := range expected {
				if v.code() == code {
					return v, nil
				}
			}
			// wrong message received
			t.Logf("received message '%v' of type '%v'\n",
				v, reflect.ValueOf(v).Type())
		}
	}
}

// private variable for mantaining engine reuse in test
// use TestEngine instead of this
var testEngine *Engine

// Engine for test reuse.
//
// Unless the test runner is passed the -no-engine-reuse flag, this will keep
// reusing the same engine.
func NewTestEngine(t *testing.T) *Engine {
	if testEngine == nil {
		opts := NewEngineOptions{Gateway: getGatewayURL()}
		if os.Getenv("CI") != "" || os.Getenv("IB_ENGINE_DUMP") != "" {
			opts.DumpConversation = true
		}
		engine, err := NewEngine(opts)

		if err != nil {
			t.Fatalf("cannot connect engine: %s", err)
		}

		if *noEngineReuse {
			t.Log("created new engine, no reuse")
			return engine
		}
		t.Log("created engine for reuse")
		testEngine = engine
		return engine
	}

	if testEngine.State() != EngineReady {
		t.Fatalf("engine %s not ready (did a prior test Stop() rather than ConditionalStop() ?)", testEngine.ConnectionInfo())
	}

	t.Logf("reusing engine %s; state: %v", testEngine.ConnectionInfo(), testEngine.State())
	return testEngine
}

// ConditionalStop will actually do a stop only if the flag -no-engine-reuse is active
func (e *Engine) ConditionalStop(t *testing.T) {
	if *noEngineReuse {
		t.Log("no engine reuse, stopping engine")
		e.Stop()
		t.Logf("engine state: %d", e.State())
	}
}

func TestConnect(t *testing.T) {
	opts := NewEngineOptions{Gateway: getGatewayURL()}
	if os.Getenv("CI") != "" || os.Getenv("IB_ENGINE_DUMP") != "" {
		opts.DumpConversation = true
	}
	engine, err := NewEngine(opts)

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	defer engine.Stop()

	if engine.State() != EngineReady {
		t.Fatalf("engine is not ready")
	}

	if engine.serverTime.IsZero() {
		t.Fatalf("server time not provided")
	}

	var states = make(chan EngineState)
	engine.SubscribeState(states)

	// stop the engine in 100 ms
	go func() {
		time.Sleep(100 * time.Millisecond)
		engine.Stop()
	}()

	newState := <-states

	if newState != EngineExitNormal {
		t.Fatalf("engine state change error")
	}

	err = engine.FatalError()
	if err != nil {
		t.Fatalf("engine reported an error: %v", err)
	}
}

func logreply(t *testing.T, reply Reply, err error) {
	if reply == nil {
		t.Logf("received reply nil")
	} else {
		t.Logf("received reply '%v' of type %v", reply, reflect.ValueOf(reply).Type())
	}
	if err != nil {
		t.Logf(" (error: '%v')", err)
	}
	t.Logf("\n")
}
