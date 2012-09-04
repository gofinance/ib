package sink

import (
	"github.com/wagerlabs/go.trade/engine"
	"testing"
	"time"
)

type symbol struct {
	id   int64
	name string
	e    *engine.Handle
	data *engine.ContractData
}

func (self *symbol) Id() int64 {
	return self.id
}

func (self *symbol) Start(e *engine.Handle) error {
	self.e = e
	req := &engine.RequestContractData{
		Symbol:       "SX7E",
		SecurityType: "IND",
	}
	req.SetId(self.id)

	if err := e.Send(req); err != nil {
		return err
	}

	return nil
}

func (self *symbol) Stop() error {
	req := &engine.CancelMarketData{}
	req.SetId(self.id)
	return self.e.Send(req)
}

func (self *symbol) Update(v engine.Reply) bool {
	switch v.(type) {
	case *engine.ContractDataEnd:
		return true
	case *engine.ContractData:
		self.data = v.(*engine.ContractData)
	}

	return false
}

func (self *symbol) Unique() string {
	return self.name
}

func TestCollection(t *testing.T) {
	e, err := engine.Make(200)

	if err != nil {
		t.Fatalf("cannot connect to engine: %s", err)
	}

	sym := &symbol{
		id:   e.NextRequestId(),
		name: "sx7e",
	}
	ch := make(chan bool)
	col := Make(e)
	col.Notify(ch)
	col.Add(sym)
	col.StartUpdate()

	select {
	case <-time.After(15 * time.Second):
		t.Fatalf("did not receive collection ready notification")
	case <-ch:
	}

	if sym.data == nil {
		t.Fatalf("no update received")
	}

	if sym1, ok := col.Lookup(sym.name); !ok {
		t.Fatalf("symbol not found in collection")
		if sym1 != sym {
			t.Fatalf("symbol retrieved from collection does not match")
		}
	}

	syms := col.Items()

	if len(syms) != 1 {
		t.Fatalf("expected 1 item in collection")
	}

	if syms[0] != sym {
		t.Fatalf("symbol in collection items does not match")
	}

	col.Cleanup()
}
