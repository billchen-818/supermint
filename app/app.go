package app

import (
	abci "github.com/vbhp/supermint/abci/types"
)

const ProtocolVersion = uint64(2)

type Application struct {
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) Info(info abci.RequestInfo) abci.ResponseInfo {
	return abci.ResponseInfo{}
}

func (a *Application) Query(query abci.RequestQuery) abci.ResponseQuery {
	return abci.ResponseQuery{}
}

func (a *Application) CheckTx(tx abci.RequestCheckTx) abci.ResponseCheckTx {
	return abci.ResponseCheckTx{}
}

func (a *Application) InitChain(chain abci.RequestInitChain) abci.ResponseInitChain {
	return abci.ResponseInitChain{}
}

func (a *Application) BeginBlock(block abci.RequestBeginBlock) abci.ResponseBeginBlock {

	return abci.ResponseBeginBlock{}
}

func (a *Application) DeliverTx(tx abci.RequestDeliverTx) abci.ResponseDeliverTx {
	return abci.ResponseDeliverTx{}
}

func (a *Application) EndBlock(block abci.RequestEndBlock) abci.ResponseEndBlock {
	return abci.ResponseEndBlock{}
}

func (a *Application) Commit() abci.ResponseCommit {
	return abci.ResponseCommit{}
}

func (a *Application) ListSnapshots(snapshots abci.RequestListSnapshots) abci.ResponseListSnapshots {
	return abci.ResponseListSnapshots{}
}

func (a *Application) OfferSnapshot(snapshot abci.RequestOfferSnapshot) abci.ResponseOfferSnapshot {
	return abci.ResponseOfferSnapshot{}
}

func (a *Application) LoadSnapshotChunk(chunk abci.RequestLoadSnapshotChunk) abci.ResponseLoadSnapshotChunk {
	return abci.ResponseLoadSnapshotChunk{}
}

func (a *Application) ApplySnapshotChunk(chunk abci.RequestApplySnapshotChunk) abci.ResponseApplySnapshotChunk {
	return abci.ResponseApplySnapshotChunk{}
}

var _ abci.Application = (*Application)(nil)
