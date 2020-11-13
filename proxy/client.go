package proxy

import (
	abcicli "github.com/vbhp/supermint/abci/client"
	"github.com/vbhp/supermint/abci/types"
	tmsync "github.com/vbhp/supermint/libs/sync"
)

// ClientCreator creates new ABCI clients.
type ClientCreator interface {
	// NewABCIClient returns a new ABCI client.
	NewABCIClient() (abcicli.Client, error)
}

//----------------------------------------------------
// local proxy uses a mutex on an in-proc app

type localClientCreator struct {
	mtx *tmsync.Mutex
	app types.Application
}

// NewLocalClientCreator returns a ClientCreator for the given app,
// which will be running locally.
func NewLocalClientCreator(app types.Application) ClientCreator {
	return &localClientCreator{
		mtx: new(tmsync.Mutex),
		app: app,
	}
}

func (l *localClientCreator) NewABCIClient() (abcicli.Client, error) {
	return abcicli.NewLocalClient(l.mtx, l.app), nil
}
