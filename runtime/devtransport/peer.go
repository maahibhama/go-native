package devtransport

import (
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
)

type Peer struct {
	connection io.ReadWriteCloser
	generation uint64
	writeMu    sync.Mutex
	pendingMu  sync.Mutex
	pending    map[uint32]chan Frame
	next       atomic.Uint32
	incoming   chan Frame
	done       chan struct{}
	closeOnce  sync.Once
}

func NewPeer(connection io.ReadWriteCloser, generation uint64) *Peer {
	peer := &Peer{connection: connection, generation: generation, pending: make(map[uint32]chan Frame), incoming: make(chan Frame, 32), done: make(chan struct{})}
	go peer.readLoop()
	return peer
}

func (p *Peer) Incoming() <-chan Frame { return p.incoming }
func (p *Peer) Done() <-chan struct{}  { return p.done }

func (p *Peer) Send(kind Kind, payload []byte) error {
	p.writeMu.Lock()
	defer p.writeMu.Unlock()
	return WriteFrame(p.connection, Frame{Kind: kind, Generation: p.generation, Payload: append([]byte(nil), payload...)})
}

func (p *Peer) Request(ctx context.Context, kind Kind, payload []byte) (Frame, error) {
	id := p.next.Add(1)
	response := make(chan Frame, 1)
	p.pendingMu.Lock()
	p.pending[id] = response
	p.pendingMu.Unlock()
	p.writeMu.Lock()
	err := WriteFrame(p.connection, Frame{Kind: kind, Generation: p.generation, RequestID: id, Payload: append([]byte(nil), payload...)})
	p.writeMu.Unlock()
	if err != nil {
		p.removePending(id)
		return Frame{}, err
	}
	select {
	case frame := <-response:
		return frame, nil
	case <-ctx.Done():
		p.removePending(id)
		return Frame{}, ctx.Err()
	case <-p.done:
		return Frame{}, errors.New("dev transport: peer closed")
	}
}

func (p *Peer) Close() error {
	var err error
	p.closeOnce.Do(func() { err = p.connection.Close(); close(p.done) })
	return err
}

func (p *Peer) readLoop() {
	defer p.Close()
	defer close(p.incoming)
	for {
		frame, err := ReadFrame(p.connection)
		if err != nil {
			return
		}
		if frame.Generation != p.generation {
			continue
		}
		if frame.RequestID != 0 {
			p.pendingMu.Lock()
			response := p.pending[frame.RequestID]
			delete(p.pending, frame.RequestID)
			p.pendingMu.Unlock()
			if response != nil {
				response <- frame
				continue
			}
		}
		select {
		case p.incoming <- frame:
		case <-p.done:
			return
		}
	}
}

func (p *Peer) removePending(id uint32) {
	p.pendingMu.Lock()
	delete(p.pending, id)
	p.pendingMu.Unlock()
}
