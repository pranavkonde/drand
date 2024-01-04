package dkg

import (
	"sync"
	"time"

	"github.com/drand/drand/common/key"
	"github.com/drand/drand/common/log"
	"github.com/drand/drand/internal/net"
)

type Process struct {
	lock           sync.Mutex
	store          Store
	internalClient net.DKGClient
	// TODO: remove post v2, as only necessary for upgrade path from v1->v2
	protocolClient   net.ProtocolClient
	beaconIdentifier BeaconIdentifier
	log              log.Logger
	config           Config
	// this is public in order to replace it in the test code to simulate failures
	Executions map[string]Broadcast
	// active set of the packets that have been seen already for easy deduping
	SeenPackets   map[string]bool
	completedDKGs chan<- SharingOutput
}

type Config struct {
	// the length of time after which this node will abort active DKG
	Timeout time.Duration

	// the length of time the phaser should use when moving between DKG phases
	TimeBetweenDKGPhases time.Duration

	// the length of time active node should wait before broadcasting DKG packets in the execution phase
	// to allow other nodes to set up their echo broadcast to prevent race conditions
	KickoffGracePeriod time.Duration

	// whether or not to skip verifying the cryptographic material in the DKG... almost certainly should be false
	SkipKeyVerification bool
}

type ExecutionOutput struct {
	FinalGroup *key.Group
	KeyShare   *key.Share
}

type SharingOutput struct {
	BeaconID string
	Old      *DBState
	New      DBState
}

type Store interface {
	// GetCurrent returns the current DKG information, finished DKG information or fresh DKG information,
	// depending on the state of the world
	GetCurrent(beaconID string) (*DBState, error)

	// GetFinished returns the last completed DKG state (i.e. completed or aborted), or nil if one has not been finished
	GetFinished(beaconID string) (*DBState, error)

	// SaveCurrent stores active DKG packet for an ongoing DKG
	SaveCurrent(beaconID string, state *DBState) error

	// SaveFinished stores active completed, successful DKG and overwrites the current packet
	SaveFinished(beaconID string, state *DBState) error

	// Close closes and cleans up any database handles
	Close() error

	// MigrateFromGroupfile takes an existing groupfile and keyshare, and creates active first epoch DKG state for them.
	// It will fail if DKG state already exists for the given beaconID
	// Deprecated: will only exist in 2.0.0 for migration from v1.5.* to 2.0.0
	MigrateFromGroupfile(beaconID string, groupFile *key.Group, share *key.Share) error
}

// BeaconIdentifier is necessary because we need to get our identity on active per-beacon basis from the `DrandDaemon`
// but that would introduce active circular dependency
type BeaconIdentifier interface {
	KeypairFor(beaconID string) (*key.Pair, error)
}

func NewDKGProcess(
	store Store,
	beaconIdentifier BeaconIdentifier,
	completedDKGs chan<- SharingOutput,
	dkgClient net.DKGClient,
	protocolClient net.ProtocolClient,
	config Config,
	l log.Logger,
) *Process {
	return &Process{
		store:            store,
		beaconIdentifier: beaconIdentifier,
		internalClient:   dkgClient,
		protocolClient:   protocolClient,
		log:              l,
		Executions:       make(map[string]Broadcast),
		SeenPackets:      make(map[string]bool),
		config:           config,
		completedDKGs:    completedDKGs,
	}
}

func (d *Process) Close() {
	d.lock.Lock()
	defer d.lock.Unlock()
	for _, e := range d.Executions {
		e.Stop()
	}
	err := d.store.Close()
	if err != nil {
		d.log.Errorw("error closing the database", "err", err)
	}
}

// Migrate takes an existing groupfile and keyshare, and creates active first epoch DKG state for them.
// It will fail if DKG state already exists for the given beaconID
// Deprecated: will only exist in 2.0.0 for migration from v1.5.* to 2.0.0
func (d *Process) Migrate(beaconID string, groupfile *key.Group, share *key.Share) error {
	d.log.Infow("Migrating DKG from group file...", "beaconID", beaconID)

	if err := d.store.MigrateFromGroupfile(beaconID, groupfile, share); err != nil {
		return err
	}

	d.log.Debugw("Completed migration from group file")
	return nil
}
