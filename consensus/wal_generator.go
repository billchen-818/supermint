package consensus

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/vbhp/supermint/app"
	"io"
	"testing"
	"time"

	db "github.com/tendermint/tm-db"

	cfg "github.com/vbhp/supermint/config"
	"github.com/vbhp/supermint/libs/log"
	tmrand "github.com/vbhp/supermint/libs/rand"
	"github.com/vbhp/supermint/privval"
	"github.com/vbhp/supermint/proxy"
	sm "github.com/vbhp/supermint/state"
	"github.com/vbhp/supermint/store"
	"github.com/vbhp/supermint/types"
)

// WALGenerateNBlocks generates a consensus WAL. It does this by spinning up a
// stripped down version of node (proxy app, event bus, consensus state) with a
// persistent kvstore application and special consensus wal instance
// (byteBufferWAL) and waits until numBlocks are created.
// If the node fails to produce given numBlocks, it returns an error.
func WALGenerateNBlocks(t *testing.T, wr io.Writer, numBlocks int) (err error) {
	config := getConfig(t)

	//ap := kvstore.NewPersistentKVStoreApplication(filepath.Join(config.DBDir(), "wal_generator"))

	logger := log.TestingLogger().With("wal_generator", "wal_generator")
	logger.Info("generating WAL (last height msg excluded)", "numBlocks", numBlocks)

	ap := app.NewApplication()
	// COPY PASTE FROM node.go WITH A FEW MODIFICATIONS
	// NOTE: we can't import node package because of circular dependency.
	// NOTE: we don't do handshake so need to set state.Version.Consensus.App directly.
	privValidatorKeyFile := config.PrivValidatorKeyFile()
	privValidatorStateFile := config.PrivValidatorStateFile()
	privValidator, err := privval.LoadOrGenFilePV(privValidatorKeyFile, privValidatorStateFile)
	if err != nil {
		return err
	}
	genDoc, err := types.GenesisDocFromFile(config.GenesisFile())
	if err != nil {
		return fmt.Errorf("failed to read genesis file: %w", err)
	}
	blockStoreDB := db.NewMemDB()
	stateDB := blockStoreDB
	stateStore := sm.NewStore(stateDB)
	state, err := sm.MakeGenesisState(genDoc)
	if err != nil {
		return fmt.Errorf("failed to make genesis state: %w", err)
	}
	state.Version.Consensus.App = app.ProtocolVersion
	if err = stateStore.Save(state); err != nil {
		t.Error(err)
	}

	blockStore := store.NewBlockStore(blockStoreDB)

	proxyApp := proxy.NewAppConns(proxy.NewLocalClientCreator(ap))
	proxyApp.SetLogger(logger.With("module", "proxy"))
	if err := proxyApp.Start(); err != nil {
		return fmt.Errorf("failed to start proxy app connections: %w", err)
	}
	t.Cleanup(func() {
		if err := proxyApp.Stop(); err != nil {
			t.Error(err)
		}
	})

	eventBus := types.NewEventBus()
	eventBus.SetLogger(logger.With("module", "events"))
	if err := eventBus.Start(); err != nil {
		return fmt.Errorf("failed to start event bus: %w", err)
	}
	t.Cleanup(func() {
		if err := eventBus.Stop(); err != nil {
			t.Error(err)
		}
	})
	mempool := emptyMempool{}
	evpool := sm.EmptyEvidencePool{}
	blockExec := sm.NewBlockExecutor(stateStore, log.TestingLogger(), proxyApp.Consensus(), mempool, evpool)
	consensusState := NewState(config.Consensus, state.Copy(), blockExec, blockStore, mempool, evpool)
	consensusState.SetLogger(logger)
	consensusState.SetEventBus(eventBus)
	if privValidator != nil {
		consensusState.SetPrivValidator(privValidator)
	}
	// END OF COPY PASTE

	// set consensus wal to buffered WAL, which will write all incoming msgs to buffer
	numBlocksWritten := make(chan struct{})
	wal := newByteBufferWAL(logger, NewWALEncoder(wr), int64(numBlocks), numBlocksWritten)
	// see wal.go#103
	if err := wal.Write(EndHeightMessage{0}); err != nil {
		t.Error(err)
	}

	consensusState.wal = wal

	if err := consensusState.Start(); err != nil {
		return fmt.Errorf("failed to start consensus state: %w", err)
	}

	select {
	case <-numBlocksWritten:
		if err := consensusState.Stop(); err != nil {
			t.Error(err)
		}
		return nil
	case <-time.After(1 * time.Minute):
		if err := consensusState.Stop(); err != nil {
			t.Error(err)
		}
		return fmt.Errorf("waited too long for supermint to produce %d blocks (grep logs for `wal_generator`)", numBlocks)
	}
}

// WALWithNBlocks returns a WAL content with numBlocks.
func WALWithNBlocks(t *testing.T, numBlocks int) (data []byte, err error) {
	var b bytes.Buffer
	wr := bufio.NewWriter(&b)

	if err := WALGenerateNBlocks(t, wr, numBlocks); err != nil {
		return []byte{}, err
	}

	wr.Flush()
	return b.Bytes(), nil
}

func randPort() int {
	// returns between base and base + spread
	base, spread := 20000, 20000
	return base + tmrand.Intn(spread)
}

func makeAddrs() (string, string, string) {
	start := randPort()
	return fmt.Sprintf("tcp://127.0.0.1:%d", start),
		fmt.Sprintf("tcp://127.0.0.1:%d", start+1),
		fmt.Sprintf("tcp://127.0.0.1:%d", start+2)
}

// getConfig returns a config for test cases
func getConfig(t *testing.T) *cfg.Config {
	c := cfg.ResetTestRoot(t.Name())

	// and we use random ports to run in parallel
	tm, rpc, grpc := makeAddrs()
	c.P2P.ListenAddress = tm
	c.RPC.ListenAddress = rpc
	c.RPC.GRPCListenAddress = grpc
	return c
}

// byteBufferWAL is a WAL which writes all msgs to a byte buffer. Writing stops
// when the heightToStop is reached. Client will be notified via
// signalWhenStopsTo channel.
type byteBufferWAL struct {
	enc               *WALEncoder
	stopped           bool
	heightToStop      int64
	signalWhenStopsTo chan<- struct{}

	logger log.Logger
}

// needed for determinism
var fixedTime, _ = time.Parse(time.RFC3339, "2017-01-02T15:04:05Z")

func newByteBufferWAL(logger log.Logger, enc *WALEncoder, nBlocks int64, signalStop chan<- struct{}) *byteBufferWAL {
	return &byteBufferWAL{
		enc:               enc,
		heightToStop:      nBlocks,
		signalWhenStopsTo: signalStop,
		logger:            logger,
	}
}

// Save writes message to the internal buffer except when heightToStop is
// reached, in which case it will signal the caller via signalWhenStopsTo and
// skip writing.
func (w *byteBufferWAL) Write(m WALMessage) error {
	if w.stopped {
		w.logger.Debug("WAL already stopped. Not writing message", "msg", m)
		return nil
	}

	if endMsg, ok := m.(EndHeightMessage); ok {
		w.logger.Debug("WAL write end height message", "height", endMsg.Height, "stopHeight", w.heightToStop)
		if endMsg.Height == w.heightToStop {
			w.logger.Debug("Stopping WAL at height", "height", endMsg.Height)
			w.signalWhenStopsTo <- struct{}{}
			w.stopped = true
			return nil
		}
	}

	w.logger.Debug("WAL Write Message", "msg", m)
	err := w.enc.Encode(&TimedWALMessage{fixedTime, m})
	if err != nil {
		panic(fmt.Sprintf("failed to encode the msg %v", m))
	}

	return nil
}

func (w *byteBufferWAL) WriteSync(m WALMessage) error {
	return w.Write(m)
}

func (w *byteBufferWAL) FlushAndSync() error { return nil }

func (w *byteBufferWAL) SearchForEndHeight(
	height int64,
	options *WALSearchOptions) (rd io.ReadCloser, found bool, err error) {
	return nil, false, nil
}

func (w *byteBufferWAL) Start() error { return nil }
func (w *byteBufferWAL) Stop() error  { return nil }
func (w *byteBufferWAL) Wait()        {}
