package symbols

import (
	"bufio"
	"bytes"
	"github.com/wagerlabs/go.trade/collection"
	"github.com/wagerlabs/go.trade/engine"
	"io"
	"os"
	"strings"
)

type Collection struct {
	e       *engine.Handle
	symbols *collection.Items
}

type Symbol struct {
	id       int64
	Name     string
	SecType  string
	Exchange string
	Currency string
	Data     []*engine.ContractData
	e        *engine.Handle
	fn       stateFn
}

type stateFn func(*Symbol) (stateFn, error)

// Make reads a list of symbols from a text file 
// and returns a contract description for each symbol
func Make(e *engine.Handle, fname string) (*Collection, error) {
	lines, err := readLines(fname)

	if err != nil {
		return nil, err
	}

	self := &Collection{
		e:       e,
		symbols: collection.Make(e),
	}

	for _, line := range lines {
		if line[0] == '#' {
			continue
		}
		var name, sectype, exchange, currency string
		v := strings.Split(line, ",")

		switch len(v) {
		case 1:
			name = v[0]
		case 2:
			name = v[0]
			sectype = v[1]
		case 3:
			name = v[0]
			sectype = v[1]
			exchange = v[2]
		case 4:
			name = v[0]
			sectype = v[1]
			exchange = v[2]
			currency = v[3]
		}

		sym := &Symbol{
			e:        e,
			id:       e.NextRequestId(),
			Name:     name,
			SecType:  sectype,
			Exchange: exchange,
			Currency: currency,
			fn:       requestStock,
			Data:     make([]*engine.ContractData, 0),
		}
		self.symbols.Add(sym)
	}

	return self, nil
}

func (self *Collection) Symbols() []*Symbol {
	src := self.symbols.Items()
	n := len(src)
	dst := make([]*Symbol, n)
	for ix, pos := range src {
		dst[ix] = pos.(*Symbol)
	}
	return dst
}

func (self *Collection) Notify(c chan bool) {
	self.symbols.Notify(c)
}

func (self *Collection) StartUpdate() {
	self.symbols.StartUpdate()
}

func (self *Symbol) Id() int64 {
	return self.id
}

func (self *Symbol) Start(e *engine.Handle) (int64, error) {
	return self.step()
}

func (self *Symbol) Stop() error {
	req := &engine.CancelMarketData{}
	req.SetId(self.id)
	return self.e.Send(req)
}

func (self *Symbol) Update(v engine.Reply) (int64, bool) {
	switch v.(type) {
	case *engine.ContractDataEnd:
		return self.id, true
	case *engine.ContractData:
		self.Data = append(self.Data, v.(*engine.ContractData))
	case *engine.ErrorMessage:
		self.step()
	}

	return self.id, false
}

func (self *Symbol) Unique() string {
	return self.Name
}

// state machine

func (self *Symbol) step() (int64, error) {
	fn, err := self.fn(self)

	if err == nil {
		self.fn = fn
		return self.id, nil
	}

	return 0, err
}

func requestStock(self *Symbol) (stateFn, error) {
	return requestIndex, self.request("STK", "SMART")
}

func requestIndex(self *Symbol) (stateFn, error) {
	return requestFuture, self.request("IND", "")
}

func requestFuture(self *Symbol) (stateFn, error) {
	return nil, self.request("FUT", "")
}

func (self *Symbol) request(sectype string, exchange string) error {
	self.id = self.e.NextRequestId()
	if self.SecType != "" {
		sectype = self.SecType
	}
	if self.Exchange != "" {
		exchange = self.Exchange
	}
	req := &engine.RequestContractData{
		Symbol:       self.Name,
		SecurityType: sectype,
		Exchange:     exchange,
		Currency:     self.Currency,
	}
	req.SetId(self.id)

	if err := self.e.Send(req); err != nil {
		return err
	}

	return nil
}

func readLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)

	if file, err = os.Open(path); err != nil {
		return
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))

	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			s := buffer.String()
			s = strings.Trim(s, " \n\r")
			lines = append(lines, s)
			buffer.Reset()
		}
	}

	if err == io.EOF {
		err = nil
	}

	return
}
