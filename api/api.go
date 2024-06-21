package api

import (
	"fmt"
	r "github.com/hashicorp/raft"
	"github.com/karlhuang95/praetor/fsm"
	"github.com/karlhuang95/praetor/raft"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var (
	isLeader int64
)

type HttpServer struct {
	ctx *r.Raft
	fsm *fsm.Fsm
}

func (h HttpServer) Set(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt64(&isLeader) == 0 {
		fmt.Fprintf(w, "not leader")
		return
	}
	vars := r.URL.Query()
	key := vars.Get("key")
	value := vars.Get("value")
	if key == "" || value == "" {
		fmt.Fprintf(w, "error key or value")
		return
	}

	data := "set" + "," + key + "," + value

	future := h.ctx.Apply([]byte(data), 5*time.Second)
	if err := future.Error(); err != nil {
		fmt.Fprintf(w, "error:"+err.Error())
		return
	}
	fmt.Fprintf(w, "ok")
	return
}

func (h HttpServer) Get(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	key := vars.Get("key")
	if key == "" {
		fmt.Fprintf(w, "error key")
		return
	}
	value := h.fsm.DataBase.Get(key)
	fmt.Fprintf(w, value)
	return
}

func (h HttpServer) Del(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt64(&isLeader) == 0 {
		fmt.Fprintf(w, "not leader")
		return
	}
	vars := r.URL.Query()
	key := vars.Get("key")
	if key == "" {
		fmt.Fprintf(w, "error key or value")
		return
	}

	data := "del" + "," + key

	future := h.ctx.Apply([]byte(data), 5*time.Second)
	if err := future.Error(); err != nil {
		fmt.Fprintf(w, "error:"+err.Error())
		return
	}
	fmt.Fprintf(w, "ok")
	return
}

func (h HttpServer) State(w http.ResponseWriter, r *http.Request) {
	state := h.ctx.State()
	fmt.Fprintf(w, state.String())

	return

}

func Start(httpAddr, raftAddr, raftId, rafCluster string) {
	if httpAddr == "" || raftAddr == "" || raftId == "" || rafCluster == "" {
		fmt.Println("config error")
		os.Exit(1)
		return
	}
	raftDir := "node/raft_" + raftId
	_, err := os.Stat(raftDir)
	if err != nil {
		os.MkdirAll(raftDir, 0700)
	}

	// 初始化raft
	newRaft, fm, err := raft.NewRaft(raftAddr, raftId, raftDir)
	if err != nil {
		fmt.Println("NewRaft error", err)
		os.Exit(1)
		return
	}

	// 启动
	raft.Bootstrap(newRaft, raftId, raftAddr, rafCluster)

	// 监听leader变化（使用此方法无法保证强一致性读，仅做leader变化过程观察）
	go func() {
		for leader := range newRaft.LeaderCh() {
			if leader {
				atomic.StoreInt64(&isLeader, 1)
			} else {
				atomic.StoreInt64(&isLeader, 0)
			}

		}
	}()

	// 启动http server
	httpServer := HttpServer{
		ctx: newRaft,
		fsm: fm,
	}

	// 接口
	http.HandleFunc("/set", httpServer.Set)
	http.HandleFunc("/get", httpServer.Get)
	http.HandleFunc("/del", httpServer.Del)
	http.HandleFunc("/state", httpServer.State)

	http.ListenAndServe(httpAddr, nil)

	// 关闭raft
	shutdownFuture := newRaft.Shutdown()
	if err := shutdownFuture.Error(); err != nil {
		fmt.Printf("shutdown raft error:%v \n", err)
	}

	// 退出http server
	fmt.Println("shutdown praetor http server")
}
