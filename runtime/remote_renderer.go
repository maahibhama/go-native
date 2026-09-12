package runtime

import (
	"errors"

	"github.com/go-native/go-native/runtime/devtransport"
)

// RemoteRenderer sends the unchanged production MutationBatch format through
// a development transport peer connected to a persistent native shell.
type RemoteRenderer struct{ Peer *devtransport.Peer }

func (renderer *RemoteRenderer) Apply(batch MutationBatch) error {
	if renderer == nil || renderer.Peer == nil {
		return errors.New("remote renderer: disconnected")
	}
	payload, err := batch.MarshalBinary()
	if err != nil {
		return err
	}
	return renderer.Peer.Send(devtransport.KindMutationBatch, payload)
}

func (renderer *RemoteRenderer) ResetTree() error {
	if renderer == nil || renderer.Peer == nil {
		return errors.New("remote renderer: disconnected")
	}
	return renderer.Peer.Send(devtransport.KindResetTree, nil)
}
