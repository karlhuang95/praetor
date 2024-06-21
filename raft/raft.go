package raft

import (
	"fmt"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/karlhuang95/praetor/fsm"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func NewRaft(raftAddr, raftId, raftDir string) (*raft.Raft, *fsm.Fsm, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(raftId)
	config.HeartbeatTimeout = 1000 * time.Millisecond
	config.ElectionTimeout = 1000 * time.Millisecond
	config.CommitTimeout = 1000 * time.Millisecond
	addr, err := net.ResolveTCPAddr("tcp", raftAddr)
	if err != nil {
		return nil, nil, err
	}
	transport, err := raft.NewTCPTransport(raftAddr, addr, 2, 5*time.Second, os.Stderr)
	if err != nil {
		return nil, nil, err
	}
	snapstore, err := raft.NewFileSnapshotStore(raftDir, 5, os.Stderr)
	if err != nil {
		return nil, nil, err
	}
	logName := fmt.Sprintf("node_%v_log.db", raftId)
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, logName))
	if err != nil {
		return nil, nil, err
	}
	stableName := fmt.Sprintf("node_%v_stable.db", raftId)
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, stableName))
	if err != nil {
		return nil, nil, err
	}
	fm := fsm.NewFsm()
	rf, err := raft.NewRaft(config, fm, logStore, stableStore, snapstore, transport)
	if err != nil {
		return nil, nil, err
	}
	return rf, fm, nil
}

func Bootstrap(rf *raft.Raft, raftId, raftAddr, raftCluster string) {
	server := rf.GetConfiguration().Configuration().Servers
	if len(server) > 0 {
		return
	}
	peerArray := strings.Split(raftCluster, ",")
	if len(peerArray) == 0 {
		return
	}
	var configuration raft.Configuration
	for _, peerInfo := range peerArray {
		peer := strings.Split(peerInfo, "/")
		id := peer[0]
		addr := peer[1]
		server := raft.Server{
			ID:      raft.ServerID(id),
			Address: raft.ServerAddress(addr),
		}
		configuration.Servers = append(configuration.Servers, server)
	}
	rf.BootstrapCluster(configuration)
	return
}
