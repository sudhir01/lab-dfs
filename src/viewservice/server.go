package viewservice

import "net"
import "net/rpc"
import "log"
import "time"
import "sync"
import "fmt"
import "os"

type ViewServer struct {
  mu          sync.Mutex
  l           net.Listener
  dead        bool
  me          string
  pingTimes   map[string] time.Time
  currentView View //TODO - add logic for current view
}

/**
* Keep track of the fact that a server has acknowledged the current view
*/
func (vs *ViewServer) updatePrimaryAck(server string, viewnum uint) {
    if vs.currentView.Primary == server {
        if vs.currentView.Viewnum == viewnum {
            vs.currentView.PrimaryAck = true
        } else {
            vs.currentView.PrimaryAck = false
        }
    }
}

//
// server Ping RPC handler.
//
func (vs *ViewServer) Ping(args *PingArgs, reply *PingReply) error {
    vs.mu.Lock()
    defer vs.mu.Unlock()

    viewnum              := args.Viewnum
    server               := args.Me
    vs.pingTimes[server] = time.Now()

    vs.updatePrimaryAck(server, viewnum)

    return nil
}

// 
// server Get() RPC handler.
//
func (vs *ViewServer) Get(args *GetArgs, reply *GetReply) error {

  // Your code here.

  return nil
}


//
// tick() is called once per PingInterval; it should notice
// if servers have died or recovered, and change the view
// accordingly.
//
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

func StartServer(me string) *ViewServer {
  vs := new(ViewServer)
  vs.me          = me
  vs.pingTimes   = map[string] time.Time{}
  vs.currentView = View{0, "", "", false}
  //vs.currentView = View{INITIAL_VIEW, NO_SERVER, NO_SERVER, false}

  // tell net/rpc about our RPC server and handlers.
  rpcs := rpc.NewServer()
  rpcs.Register(vs)

  // prepare to receive connections from clients.
  // change "unix" to "tcp" to use over a network.
  os.Remove(vs.me) // only needed for "unix"
  l, e := net.Listen("unix", vs.me);
  if e != nil {
    log.Fatal("listen error: ", e);
  }
  vs.l = l

  // please don't change any of the following code,
  // or do anything to subvert it.

  // create a thread to accept RPC connections from clients.
  go func() {
    for vs.dead == false {
      conn, err := vs.l.Accept()
      if err == nil && vs.dead == false {
        go rpcs.ServeConn(conn)
      } else if err == nil {
        conn.Close()
      }
      if err != nil && vs.dead == false {
        fmt.Printf("ViewServer(%v) accept: %v\n", me, err.Error())
        vs.Kill()
      }
    }
  }()

  // create a thread to call tick() periodically.
  go func() {
    for vs.dead == false {
      vs.tick()
      time.Sleep(PingInterval)
    }
  }()

  return vs
}
