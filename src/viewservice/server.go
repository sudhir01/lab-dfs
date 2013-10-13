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

func (this *ViewServer) Start() {
    this.registerRPCServer()
    this.openPort()
    go this.startConnectionAcceptor()
    go this.startTicker()
}

//TODO - deprecated method, update references and remove it
func StartServer(me string) *ViewServer {
	 rpc     := rpc.NewServer()
	 timer   := new(DefaultTimer)
	 tracker := NewViewTracker(timer)
	 handler := NewViewServerHandler(tracker)
	 this, _ := NewViewServer(me, rpc, tracker, handler)
	 this.Start()
	 return this
}

func (this *ViewServer) IsListening() bool {
    return this.l != nil
}

func (this *ViewServer) IsDead() bool {
    return this.dead
}

func (this *ViewServer) Name() string {
    return this.me
}

func (this *ViewServer) ListenerAddress() string {
    if this.l == nil {
        return ""
    }
    return this.l.Addr().String()
}

//
// tell the server to shut itself down.
// for testing.
// please don't change this function.
//
func (this *ViewServer) Kill() {
  this.dead = true
  this.l.Close()
}

//private methods

func (this *ViewServer) markDeadServers() {
}

func (this *ViewServer) tick() {
	 this.tracker.Tick()
}

func (this *ViewServer) openPort() {
    // prepare to receive connections from clients.
    // change "unix" to "tcp" to use over a network.
    os.Remove(this.me) // only needed for "unix"
    l, e := net.Listen("unix", this.me)
    if e != nil {
        log.Fatal("listen error: ", e)
    }
    this.l = l
}

func (this *ViewServer) registerRPCServer() {
    this.rpcServer.RegisterName("ViewServer", this.handler)
}

func (this *ViewServer) dispatch(conn net.Conn) {
	 go this.rpcServer.ServeConn(conn)
}

func (this *ViewServer) startConnectionAcceptor() {
    for this.dead == false {
        conn, err := this.l.Accept()
        if err == nil && this.dead == false {
				this.dispatch(conn)
        } else if err == nil {
            conn.Close()
        }
        if err != nil && this.dead == false {
            log.Print("ViewServer(%v) accept: %v\n", this.me, err.Error())
            this.Kill()
        }
    }
}

func (this *ViewServer) startTicker() {
    for this.dead == false {
        this.tick()
        time.Sleep(PingInterval)
    }
}
