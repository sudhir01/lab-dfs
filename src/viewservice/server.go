package viewservice

import "errors"
import "net"
import "net/rpc"
import "log"
import "time"
import "sync"
import "os"

type ViewServer struct {
  rpcServer   *rpc.Server
  mu          sync.Mutex
  l           net.Listener
  dead        bool
  me          string
  pingTimes   map[string] time.Time
  currentView View
  handler     ServerHandler
}

func (vs *ViewServer) IsListening() bool {
    return vs.l != nil
}

func (vs *ViewServer) IsDead() bool {
    return vs.dead
}

func (vs *ViewServer) Name() string {
    return vs.me
}

func (vs *ViewServer) ListenerAddress() string {
    if vs.l == nil {
        return ""
    }
    return vs.l.Addr().String()
}

func (vs *ViewServer) View() *View {
    return &vs.currentView
}

func (vs *ViewServer) PingTable() *map[string] time.Time {
    return &vs.pingTimes
}

func (vs *ViewServer) hasPrimaryAck() bool {
    return (vs.currentView.Viewnum == INITIAL_VIEW) || (vs.currentView.PrimaryView == vs.currentView.Viewnum)
}

// server Ping RPC handler.
// If ping payload is 0, then the server crashed
func (vs *ViewServer) Ping(args *PingArgs, reply *PingReply) error {
    vs.mu.Lock()
    defer vs.mu.Unlock()

    viewnum              := args.Viewnum
    server               := args.Me
    vs.pingTimes[server] = time.Now()

    switch server {
    case vs.currentView.Primary:
        vs.currentView.PrimaryView = viewnum
    case vs.currentView.Backup:
        vs.currentView.BackupView  = viewnum
    }

    reply.View = vs.currentView
    return nil
}

func elapsedDeadPings(lastPing time.Time) int64 {
    now   := time.Now()
    delta := now.Sub(lastPing)
    return int64(delta/PingInterval)
    //milli := (delta.Nanoseconds()/1000)
    //pings := milli / (PingInterval.Nanoseconds()/1000)
    //return pings
}

func (vs *ViewServer) markDeadServers() {
}

// tick() is called once per PingInterval; it should notice
// if servers have died or recovered, and change the view
// accordingly.
// Periodic tasks on ping:
//    1. Mark server dead if max DeadPings have passed for PingIntervals
//    2. Update view for either (only if primary has not drifted):
//       i.  dead server or                      -> TODO
//       ii. idle server when there is no backup -> TODO
func (vs *ViewServer) tick() {

  // Your code here.
}

//
// tell the server to shut itself down.
// for testing.
// please don't change this function.
//
func (vs *ViewServer) Kill() {
  vs.dead = true
  vs.l.Close()
}

func NewViewServer(hostPort string, rpcServer *rpc.Server, handler ServerHandler) (*ViewServer, error) {
    if hostPort == "" {
        err := errors.New("hostPort cannot be empty")
        return nil, err
    }

    if rpcServer == nil {
        err := errors.New("RPC server cannot be nil")
        return nil, err
    }

	 if handler == nil {
		  err := errors.New("Server handler cannot be nil")
		  return nil, err
	 }

    vs := new(ViewServer)

    vs.me          = hostPort
    vs.pingTimes   = map[string] time.Time{}
    vs.currentView = View{INITIAL_VIEW, NO_SERVER, NO_SERVER, NO_VIEW, NO_VIEW}
    vs.rpcServer   = rpcServer
	 vs.handler     = handler
    return vs, nil
}

func (vs *ViewServer) Start() {
    vs.registerRPCServer()
    vs.openPort()
    go vs.startConnectionAcceptor()
    go vs.startTicker()
}

func (vs *ViewServer) openPort() {
    // prepare to receive connections from clients.
    // change "unix" to "tcp" to use over a network.
    os.Remove(vs.me) // only needed for "unix"
    l, e := net.Listen("unix", vs.me)
    if e != nil {
        log.Fatal("listen error: ", e)
    }
    vs.l = l
}

func (vs *ViewServer) registerRPCServer() {
	 //FIXME - registration can happen independent of view server
    vs.rpcServer.Register(vs.handler)
}

func (vs *ViewServer) dispatch(conn net.Conn) {
	 go vs.rpcServer.ServeConn(conn)
}

func (vs *ViewServer) startConnectionAcceptor() {
    for vs.dead == false {
        conn, err := vs.l.Accept()
        if err == nil && vs.dead == false {
				vs.dispatch(conn)
        } else if err == nil {
            conn.Close()
        }
        if err != nil && vs.dead == false {
            log.Print("ViewServer(%v) accept: %v\n", vs.me, err.Error())
            vs.Kill()
        }
    }
}

func (vs *ViewServer) startTicker() {
    for vs.dead == false {
        vs.tick()
        time.Sleep(PingInterval)
    }
}

func StartServer(me string) *ViewServer {
    rpc     := rpc.NewServer()
	 handler := new(ViewServerHandler)
    vs, _ := NewViewServer(me, rpc, handler)
    vs.Start()
    return vs
}
