package viewservice

import "net"
import "net/rpc"
import "log"
import "time"
import "os"

type ViewServer struct {
  rpcServer   *rpc.Server
  l           net.Listener
  dead        bool
  me          string
  handler     ServerHandler
  tracker	  *ViewTracker
}

func (vs *ViewServer) Start() {
    vs.registerRPCServer()
    vs.openPort()
    go vs.startConnectionAcceptor()
    go vs.startTicker()
}

//TODO - deprecated method, update references and remove it
func StartServer(me string) *ViewServer {
	 rpc     := rpc.NewServer()
	 tracker := NewViewTracker()
	 handler := NewViewServerHandler(tracker)
	 vs, _ := NewViewServer(me, rpc, tracker, handler)
	 vs.Start()
	 return vs
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

//
// tell the server to shut itself down.
// for testing.
// please don't change this function.
//
func (vs *ViewServer) Kill() {
  vs.dead = true
  vs.l.Close()
}

//private methods

func (vs *ViewServer) markDeadServers() {
}

func (vs *ViewServer) tick() {
	 vs.tracker.Tick()
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
    vs.rpcServer.RegisterName("ViewServer", vs.handler)
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

