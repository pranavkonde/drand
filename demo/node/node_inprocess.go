package node

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/kabukky/httpscerts"

	"github.com/drand/drand/chain"
	"github.com/drand/drand/client/grpc"
	"github.com/drand/drand/core"
	"github.com/drand/drand/crypto"
	"github.com/drand/drand/demo/cfg"
	"github.com/drand/drand/fs"
	"github.com/drand/drand/key"
	"github.com/drand/drand/log"
	"github.com/drand/drand/net"
	"github.com/drand/drand/protobuf/drand"
	"github.com/drand/drand/test"
)

// LocalNode ...
type LocalNode struct {
	base       string
	i          int
	period     string
	beaconID   string
	scheme     *crypto.Scheme
	logPath    string
	privAddr   string
	pubAddr    string
	ctrlAddr   string
	ctrlClient *net.ControlClient
	tls        bool
	priv       *key.Pair

	dbEngineType chain.StorageType
	pgDSN        func() string
	memDBSize    int

	log log.Logger

	daemon *core.DrandDaemon
}

func NewLocalNode(i int, bindAddr string, cfg cfg.Config) *LocalNode {
	nbase := path.Join(cfg.BasePath, fmt.Sprintf("node-%d", i))
	os.MkdirAll(nbase, 0740)
	logPath := path.Join(nbase, "log")

	// make certificates for the node.
	err := httpscerts.Generate(
		path.Join(nbase, fmt.Sprintf("server-%d.crt", i)),
		path.Join(nbase, fmt.Sprintf("server-%d.key", i)),
		bindAddr)
	if err != nil {
		return nil
	}
	l := &LocalNode{
		base:         nbase,
		i:            i,
		period:       cfg.Period,
		tls:          cfg.WithTLS,
		logPath:      logPath,
		log:          log.NewLogger(nil, log.LogDebug),
		pubAddr:      test.FreeBind(bindAddr),
		privAddr:     test.FreeBind(bindAddr),
		ctrlAddr:     test.FreeBind("localhost"),
		scheme:       cfg.Scheme,
		beaconID:     cfg.BeaconID,
		dbEngineType: cfg.DBEngineType,
		pgDSN:        cfg.PgDSN,
		memDBSize:    cfg.MemDBSize,
	}

	var priv *key.Pair
	if l.tls {
		priv, err = key.NewTLSKeyPair(l.privAddr, nil)
	} else {
		priv, err = key.NewKeyPair(l.privAddr, nil)
	}
	if err != nil {
		panic(err)
	}

	l.priv = priv
	return l
}

func (l *LocalNode) Start(certFolder string, dbEngineType chain.StorageType, pgDSN func() string, memDBSize int) error {
	if dbEngineType != "" {
		l.dbEngineType = dbEngineType
	}
	if pgDSN != nil {
		l.pgDSN = pgDSN
	}
	if memDBSize != 0 {
		l.memDBSize = memDBSize
	}

	certs, err := fs.Files(certFolder)
	if err != nil {
		return err
	}

	opts := []core.ConfigOption{
		core.WithLogLevel(log.LogDebug, false),
		core.WithConfigFolder(l.base),
		core.WithTrustedCerts(certs...),
		core.WithPublicListenAddress(l.pubAddr),
		core.WithPrivateListenAddress(l.privAddr),
		core.WithControlPort(l.ctrlAddr),
		core.WithDBStorageEngine(l.dbEngineType),
		core.WithPgDSN(l.pgDSN()),
		core.WithMemDBSize(l.memDBSize),
	}

	if l.tls {
		opts = append(opts, core.WithTLS(
			path.Join(l.base, fmt.Sprintf("server-%d.crt", l.i)),
			path.Join(l.base, fmt.Sprintf("server-%d.key", l.i))))
	} else {
		opts = append(opts, core.WithInsecure())
	}

	conf := core.NewConfig(opts...)
	ks := key.NewFileStore(conf.ConfigFolderMB(), l.beaconID)
	err = ks.SaveKeyPair(l.priv)
	if err != nil {
		return err
	}

	err = key.Save(path.Join(l.base, "public.toml"), l.priv.Public, false)
	if err != nil {
		return err
	}

	// Create and start drand daemon
	drandDaemon, err := core.NewDrandDaemon(conf)
	if err != nil {
		return fmt.Errorf("can't instantiate drand daemon %s", err)
	}

	// Load possible existing stores
	stores, err := key.NewFileStores(conf.ConfigFolderMB())
	if err != nil {
		return err
	}

	for beaconID, ks := range stores {
		bp, err := drandDaemon.InstantiateBeaconProcess(beaconID, ks)
		if err != nil {
			fmt.Printf("beacon id [%s]: can't instantiate randomness beacon. err: %s \n", beaconID, err)
			return err
		}

		freshRun, err := bp.Load()
		if err != nil {
			return err
		}

		if freshRun {
			fmt.Printf("beacon id [%s]: will run as fresh install -> expect to run DKG.\n", beaconID)
		} else {
			fmt.Printf("beacon id [%s]: will already start running randomness beacon.\n", beaconID)
			// Add beacon handler from chain hash for http server
			drandDaemon.AddBeaconHandler(beaconID, bp)

			// XXX make it configurable so that new share holder can still start if
			// nobody started.
			// drand.StartBeacon(!c.Bool(pushFlag.Name))
			catchup := true
			err = bp.StartBeacon(catchup)
			if err != nil {
				return err
			}
		}
	}

	l.daemon = drandDaemon

	return nil
}

func (l *LocalNode) PrivateAddr() string {
	return l.privAddr
}

func (l *LocalNode) CtrlAddr() string {
	return l.ctrlAddr
}

func (l *LocalNode) PublicAddr() string {
	return l.pubAddr
}

func (l *LocalNode) Index() int {
	return l.i
}

func (l *LocalNode) ctrl() *net.ControlClient {
	if l.ctrlClient != nil {
		return l.ctrlClient
	}
	cl, err := net.NewControlClient(l.ctrlAddr)
	if err != nil {
		l.log.Errorw("", "drand", "can't instantiate control client", "err", err)
		return nil
	}
	l.ctrlClient = cl
	return cl
}

func (l *LocalNode) RunDKG(nodes, thr int, timeout time.Duration, leader bool, leaderAddr string, beaconOffset int) (*key.Group, error) {
	cl := l.ctrl()
	p, err := time.ParseDuration(l.period)
	if err != nil {
		l.log.Errorw("", "drand", "dkg run failed", "err", err)
		return nil, err
	}
	var grp *drand.GroupPacket
	if leader {
		grp, err = cl.InitDKGLeader(nodes, thr, p, 0, timeout, nil, secretDKG, beaconOffset, l.beaconID)
	} else {
		leader := net.CreatePeer(leaderAddr, l.tls)
		grp, err = cl.InitDKG(leader, nil, secretDKG, l.beaconID)
	}
	if err != nil {
		l.log.Errorw("", "drand", "dkg run failed", "err", err)
		return nil, err
	}
	return key.GroupFromProto(grp, nil)
}

func (l *LocalNode) GetGroup() *key.Group {
	cl := l.ctrl()

	grp, err := cl.GroupFile(l.beaconID)
	if err != nil {
		l.log.Errorw("", "drand", "can't  get group", "err", err)
		return nil
	}
	group, err := key.GroupFromProto(grp, nil)
	if err != nil {
		l.log.Errorw("", "drand", "can't deserialize group", "err", err)
		return nil
	}
	return group
}

func (l *LocalNode) RunReshare(nodes, thr int, oldGroup string, timeout string, leader bool, leaderAddr string, beaconOffset int) *key.Group {
	cl := l.ctrl()

	t, _ := time.ParseDuration(timeout)
	var grp *drand.GroupPacket
	var err error
	if leader {
		grp, err = cl.InitReshareLeader(nodes, thr, t, 0, secretReshare, oldGroup, beaconOffset, l.beaconID)
	} else {
		leader := net.CreatePeer(leaderAddr, l.tls)
		grp, err = cl.InitReshare(leader, secretReshare, oldGroup, false, l.beaconID)
	}
	if err != nil {
		l.log.Errorw("", "drand", "reshare failed", "err", err)
		return nil
	}
	kg, _ := key.GroupFromProto(grp, nil)
	return kg
}

func (l *LocalNode) ChainInfo(group string) bool {
	cl := l.ctrl()
	ci, err := cl.ChainInfo(l.beaconID)
	if err != nil {
		l.log.Errorw("", "drand", "can't get chain-info", "err", err)
		return false
	}
	sdist := hex.EncodeToString(ci.PublicKey)
	fmt.Printf("\t- Node %s has chain-info %s\n", l.PrivateAddr(), sdist[10:14])
	return true
}

func (l *LocalNode) Ping() bool {
	cl := l.ctrl()
	if err := cl.Ping(); err != nil {
		l.log.Errorw("", "drand", "can't ping", "err", err)
		return false
	}
	return true
}

func (l *LocalNode) GetBeacon(groupPath string, round uint64) (resp *drand.PublicRandResponse, cmd string) {
	cert := ""
	if l.tls {
		cert = path.Join(l.base, fmt.Sprintf("server-%d.crt", l.i))
	}
	c, _ := grpc.New(l.privAddr, cert, cert == "", []byte(""))

	group := l.GetGroup()
	if group == nil {
		l.log.Errorw("", "drand", "can't get group")
		return
	}

	var err error
	cmd = "unused"
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	r, err := c.Get(ctx, round)
	if err != nil || r == nil {
		l.log.Errorw("", "drand", "can't get becon", "err", err)
	}
	if r == nil {
		return
	}
	resp = &drand.PublicRandResponse{
		Round:      r.Round(),
		Signature:  r.Signature(),
		Randomness: r.Randomness(),
	}
	return
}

func (l *LocalNode) WriteCertificate(p string) {
	if l.tls {
		exec.Command("cp", path.Join(l.base, fmt.Sprintf("server-%d.crt", l.i)), p).Run()
	}
}

func (l *LocalNode) WritePublic(p string) {
	key.Save(p, l.priv.Public, false)
}

func (l *LocalNode) Stop() {
	cl := l.ctrl()
	_, err := cl.Shutdown("")
	if err != nil {
		l.log.Errorw("", "drand", "failed to shutdown", "err", err)
		return
	}
	<-l.daemon.WaitExit()
}

func (l *LocalNode) PrintLog() {
	fmt.Printf("[-] Printing logs of node %s:\n", l.PrivateAddr())
	buff, err := os.ReadFile(l.logPath)
	if err != nil {
		fmt.Printf("[-] Can't read logs at %s !\n\n", l.logPath)
		return
	}

	fmt.Printf("%s\n", string(buff))
}
