package dkg

import (
	bytes2 "bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/BurntSushi/toml"

	pdkg "github.com/drand/drand/v2/protobuf/dkg"

	"github.com/drand/drand/v2/common/key"
	"github.com/drand/drand/v2/common/log"
)

// StoreFolder is the directory where dkg.toml and dkg.staged.toml are written
const StoreFolder = "dkg"
const FileName = "dkg.toml"
const StagedFileName = "dkg.staged.toml"

const DirPerm = 0755

type FileStore struct {
	baseFolder    string
	log           log.Logger
	migrationLock sync.Mutex
}

func NewDKGStore(baseFolder string, logLevel int) (*FileStore, error) {
	dkgStoreFolder := path.Join(baseFolder, StoreFolder)
	err := os.MkdirAll(dkgStoreFolder, DirPerm)
	if err != nil {
		return nil, err
	}
	return &FileStore{
		baseFolder: dkgStoreFolder,
		log:        log.New(nil, logLevel, true),
	}, nil
}

func getFromFilePath(path string) (*DBState, error) {
	t := DBStateTOML{}
	_, err := toml.DecodeFile(path, &t)
	if err != nil {
		return nil, err
	}
	state, err := t.FromTOML()
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (fs *FileStore) GetCurrent(beaconID string) (*DBState, error) {
	fs.migrationLock.Lock()
	defer fs.migrationLock.Unlock()
	f, err := getFromFilePath(path.Join(fs.baseFolder, beaconID, StagedFileName))
	if errors.Is(err, os.ErrNotExist) {
		fs.log.Debug("No DKG file found, returning new state")
		return NewFreshState(beaconID), nil
	}
	return f, err
}

func (fs *FileStore) GetFinished(beaconID string) (*DBState, error) {
	fs.migrationLock.Lock()
	defer fs.migrationLock.Unlock()
	f, err := getFromFilePath(path.Join(fs.baseFolder, beaconID, FileName))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	return f, err
}

func saveTOMLToFilePath(filepath string, state *DBState) error {
	w, err := os.Create(filepath)
	if err != nil {
		return err
	}
	t := state.TOML()
	err = toml.NewEncoder(w).Encode(&t)
	closeErr := w.Close()
	return errors.Join(closeErr, err)
}

// SaveCurrent stores a DKG packet for an ongoing DKG
func (fs *FileStore) SaveCurrent(beaconID string, state *DBState) error {
	fs.migrationLock.Lock()
	defer fs.migrationLock.Unlock()
	err := os.MkdirAll(path.Join(fs.baseFolder, beaconID), DirPerm)
	if err != nil {
		return err
	}
	return saveTOMLToFilePath(path.Join(fs.baseFolder, beaconID, StagedFileName), state)
}

// SaveFinished stores a completed, successful DKG and overwrites the current packet
func (fs *FileStore) SaveFinished(beaconID string, state *DBState) error {
	err := os.MkdirAll(path.Join(fs.baseFolder, beaconID), DirPerm)
	if err != nil {
		return err
	}
	err = saveTOMLToFilePath(path.Join(fs.baseFolder, beaconID, StagedFileName), state)
	if err != nil {
		return err
	}
	return saveTOMLToFilePath(path.Join(fs.baseFolder, beaconID, FileName), state)
}

func (fs *FileStore) Close() error {
	// Nothing to do for flat-file management
	return nil
}

func (fs *FileStore) MigrateFromGroupfile(beaconID string, groupFile *key.Group, share *key.Share) error {
	fs.log.Debug(fmt.Sprintf("Converting group file for beaconID %s ...", beaconID))
	if beaconID == "" {
		return errors.New("you must pass a beacon ID")
	}
	if groupFile == nil {
		return errors.New("you cannot migrate without passing a previous group file")
	}
	if share == nil {
		return errors.New("you cannot migrate without a previous distributed key share")
	}

	fs.migrationLock.Lock()
	defer fs.migrationLock.Unlock()
	dbState, err := GroupFileToDBState(beaconID, groupFile, share)
	if err != nil {
		return err
	}

	dkgFilePath := path.Join(fs.baseFolder, beaconID, FileName)
	_, err = os.Stat(dkgFilePath)
	if err == nil {
		return fmt.Errorf("Found existing DKG store at %s, aborting migration", dkgFilePath)
	}
	if !errors.Is(err, os.ErrNotExist) {
		fs.log.Debug(fmt.Sprintf("Unexpected error checking for DKG store %s: %q", dkgFilePath, err))
		return err
	}

	// Save the DKG into a toml file
	fs.log.Debug(fmt.Sprintf("Writing DKG file %s for for beaconID %s ...", dkgFilePath, beaconID))
	err = os.MkdirAll(path.Join(fs.baseFolder, beaconID), DirPerm)
	if err != nil {
		return err
	}
	if err = saveTOMLToFilePath(dkgFilePath, dbState); err != nil {
		return err
	}

	// Save the DKG into the StagedDKG toml file too
	stagedDkgFilePath := path.Join(fs.baseFolder, beaconID, StagedFileName)
	fs.log.Debug(fmt.Sprintf("Writing DKG file %s for for beaconID %s ...", stagedDkgFilePath, beaconID))
	return saveTOMLToFilePath(stagedDkgFilePath, dbState)
}

func encodeState(state *DBState) ([]byte, error) {
	var bytes []byte
	b := bytes2.NewBuffer(bytes)
	err := toml.NewEncoder(b).Encode(state.TOML())
	if err != nil {
		return nil, err
	}
	return b.Bytes(), err
}

func GroupFileToDBState(beaconID string, groupFile *key.Group, share *key.Share) (*DBState, error) {
	// map all the nodes from the group file into `drand.Participant`s
	participants := make([]*pdkg.Participant, len(groupFile.Nodes))

	if len(groupFile.Nodes) == 0 {
		return nil, errors.New("you cannot migrate from a group file that doesn't contain node info")
	}
	for i, node := range groupFile.Nodes {
		pk, err := node.Key.MarshalBinary()
		if err != nil {
			return nil, err
		}

		// MIGRATION PATH: the signature is `nil` here due to an incompatibility between v1 and v2 sigs over pub keys
		// the new signature will be filled in on first proposal using the new DKG
		participants[i] = &pdkg.Participant{
			Address:   node.Address(),
			Key:       pk,
			Signature: nil,
		}
	}

	// create an epoch 1 state with the 0th node as the leader
	return &DBState{
		BeaconID:      beaconID,
		Epoch:         1,
		State:         Complete,
		Threshold:     uint32(groupFile.Threshold),
		Timeout:       time.Now(),
		SchemeID:      groupFile.Scheme.Name,
		GenesisTime:   time.Unix(groupFile.GenesisTime, 0),
		GenesisSeed:   groupFile.GenesisSeed,
		CatchupPeriod: groupFile.CatchupPeriod,
		BeaconPeriod:  groupFile.Period,
		Leader:        participants[0],
		Remaining:     nil,
		Joining:       participants,
		Leaving:       nil,
		Acceptors:     participants,
		Rejectors:     nil,
		FinalGroup:    groupFile,
		KeyShare:      share,
	}, nil
}

// NukeState deletes the directory corresponding to the specified beaconID
func (fs *FileStore) NukeState(beaconID string) error {
	return os.RemoveAll(fs.baseFolder)
}
