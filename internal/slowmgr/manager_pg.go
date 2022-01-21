package slowmgr

import (
	"fmt"
	"github.com/hpcloud/tail"
	"log"
	"strings"
)

type ManagerPgOptions struct {
	Store  Store
	Parser Parser
}

// ManagerPg Postgres implementation of Manager
// TODO should accept io.Reader instead of file path
type ManagerPg struct {
	Options     ManagerPgOptions
	LogFilePath string
	started     bool
	stopChan    chan struct{}
	errChan     chan error
	queryChan   chan *Slow
}

func NewManagerPg(lfp string, opts ManagerPgOptions) *ManagerPg {
	if opts.Store == nil {
		opts.Store = NewStoreDefault(DefaultStoreOptions{})
	}
	if opts.Parser == nil {
		opts.Parser = &DefaultParser{}
	}
	return &ManagerPg{
		LogFilePath: lfp,
		Options:     opts,
		errChan:     make(chan error),
		stopChan:    make(chan struct{}),
		queryChan:   make(chan *Slow),
	}
}

// Start is non-blocking. Caller is not responsible
// for handling concurrency
func (m *ManagerPg) Start() error {
	t, err := tail.TailFile(m.LogFilePath, tail.Config{Follow: true})
	if err != nil {
		return err
	}

	m.started = true
	go func() {
		for {
			select {
			case l := <-t.Lines:
				go m.handleNewLogEntry(l.Text)
			case sq := <-m.queryChan:
				log.Println(fmt.Sprintf("Start(): slow query received. Query: %v", sq))
			case err := <-m.errChan:
				log.Println(fmt.Sprintf("Start(): handle error | error: %v", err))
			case <-m.stopChan:
				m.started = false
				break
			}
		}
	}()
	return nil
}

func (m *ManagerPg) handleNewLogEntry(str string) {
	r := strings.NewReader(str)
	sl, err := m.Options.Parser.Parse(r)
	if err != nil {
		m.errChan <- err
		return
	}
	err = m.Save(sl)
	if err != nil {
		m.errChan <- err
		return
	}
	m.queryChan <- sl
}

func (m *ManagerPg) Stop() {
	m.stopChan <- struct{}{}
}

func (m *ManagerPg) Filter(op string, pageNo int, itemSize int) []*Slow {
	return m.Options.Store.Filter(op, pageNo, itemSize)
}

func (m *ManagerPg) Save(s *Slow) error {
	return m.Options.Store.Add(s)
}

func (m *ManagerPg) Sort(pageNo int, itemSize int) []*Slow {
	return m.Options.Store.Sort(pageNo, itemSize)
}
