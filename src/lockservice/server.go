package lockservice

import "net"
import "net/rpc"
import "log"
import "sync"
import "os"
import "io"
import "time"

type LockServer struct {
    mu         sync.Mutex
    l          net.Listener
    dead       bool  // for test_test.go
    dying      bool // for test_test.go

    am_primary bool // am I the primary?
    name       string   // server's port
    backup     string   // backup's port

    // for each lock name, is it locked?
    locks      map[string]bool
    requestIds map[int64]bool
}

func (server *LockServer) callServer (remoteServer string,
                                      rpcname      string,
                                      args         interface{},
                                      reply        interface{}) bool {
    log.Printf("[debug] [%v] Server::callServer calling server(%v)\n", server.name, remoteServer)
    connection, errx := rpc.Dial("unix", remoteServer)
    if errx != nil {
        log.Printf("[error] [%v] Server::callServer calling server(%v) failed\n", server.name, remoteServer)
        return false
    }
    defer connection.Close()

    log.Printf("[debug] [%v] Server::callServer rpc(%v) on server(%v), attempting\n", server.name, rpcname, remoteServer)
    err := connection.Call(rpcname, args, reply)
    if err == nil {
        log.Printf("[debug] [%v] Server::callServer call rpc(%v) on server(%v), success\n", server.name, rpcname, remoteServer)
        return true
    }

    log.Printf("[error] [%v] Server::callServer call rpc(%v) on server(%v), failure\n", server.name, rpcname, remoteServer)
    return false
}

func (server *LockServer) onBackup(operation func() error) error {
    log.Printf("[debug] [%v] Server::onBackup is_primary %v\n", server.name, server.am_primary)
    if server.am_primary == true {
        return operation()
    }
    return nil
}

func (server *LockServer) lockBackup(args *LockArgs,
                                    reply *LockReply) error {
    log.Printf("[debug] [%v] Server::lockBackup calling backup server\n", server.name)
    ok := server.callServer(server.backup, "LockServer.Lock", args, reply)
    if ok == true {
        reply.OK = true
        //TODO - do something if backup responsds
    } else {
        reply.OK = false
        //TODO - call to backup failed, mark is as dead?
    }
    return nil
}

func (server *LockServer) unlockBackup(args  *UnlockArgs,
                                       reply *UnlockReply) error {
    log.Printf("[debug] [%v] Server::unlockBackup calling backup server\n", server.name)
    ok := server.callServer(server.backup, "LockServer.Unlock", args, reply)
    if ok == true {
        reply.OK = true
        //TODO - do something if backup responsds
    } else {
        reply.OK = false
        //TODO - call to backup failed, mark is as dead?
    }
    return nil
}

//
// server Lock RPC handler.
//
// you will have to modify this function
//
func (ls *LockServer) Lock(args  *LockArgs,
                           reply *LockReply) error {
    log.Printf("[debug] [%v] Server::Lock lock (%v) called for server, locks: (%v)\n", ls.name, args, ls.locks)
    ls.mu.Lock()
    defer ls.mu.Unlock()

    requestId           := args.RequestId
    duplicateRequest, _ := ls.requestIds[requestId]
    locked, _           := ls.locks[args.Lockname]

    log.Printf("[debug] [%v] Server::Lock current state duplicateRequest(%v) locked(%v)\n", ls.name, duplicateRequest, locked)

    if duplicateRequest && locked {
        reply.OK = true
        return nil
    }

    if locked {
        reply.OK = false
    } else {
        reply.OK = true
        ls.locks[args.Lockname]  = true
        ls.requestIds[requestId] = true
    }

    var backupReply LockReply
    log.Printf("[debug] [%v] Server::Lock calling backup server\n", ls.name)
    backupLockFn := func() error { return ls.lockBackup(args, &backupReply) }
    ls.onBackup(backupLockFn)
    return nil
}

//
// server Unlock RPC handler.
//
func (ls *LockServer) Unlock(args  *UnlockArgs,
                             reply *UnlockReply) error {
    ls.mu.Lock()
    defer ls.mu.Unlock()

    var unlockReply UnlockReply
    fn := func() error { return ls.unlockBackup(args, &unlockReply) }
    ls.onBackup(fn)
    requestId           := args.RequestId
    duplicateRequest, _ := ls.requestIds[requestId]
    locked, _           := ls.locks[args.Lockname]

    if duplicateRequest && !locked {
        reply.OK = true
        return nil
    }

    if locked {
        ls.locks[args.Lockname]  = false
        ls.requestIds[requestId] = true
        reply.OK = true
    } else {
        reply.OK = false
    }

    return nil
}

//
// tell the server to shut itself down.
// for testing.
// please don't change this.
//
func (ls *LockServer) kill() {
    ls.dead = true
    ls.l.Close()
}

//
// hack to allow test_test.go to have primary process
// an RPC but not send a reply. can't use the shutdown()
// trick b/c that causes client to immediately get an
// error and send to backup before primary does.
// please don't change anything to do with DeafConn.
//
type DeafConn struct {
    c io.ReadWriteCloser
}
func (dc DeafConn) Write(p []byte) (n   int,
                                    err error) {
    return len(p), nil
}
func (dc DeafConn) Close() error {
    return dc.c.Close()
}
func (dc DeafConn) Read(p []byte) (n   int,
                                   err error) {
    return dc.c.Read(p)
}

func StartServer(primary    string,
                 backup     string,
                 am_primary bool) *LockServer {
    log.Printf("[debug] Server::StartServer starting server primary(%v), backup(%v), am_primary(%v)\n", primary, backup, am_primary)
    ls := new(LockServer)
    if am_primary == true {
        ls.name = primary
    } else {
        ls.name = backup
    }
    ls.backup = backup
    ls.am_primary = am_primary
    ls.locks = map[string]bool{}
    ls.requestIds = map[int64]bool{}

    // Your initialization code here.

    me := ""
    if am_primary {
        me = primary
    } else {
        me = backup
    }

    // tell net/rpc about our RPC server and handlers.
    rpcs := rpc.NewServer()
    rpcs.Register(ls)

    // prepare to receive connections from clients.
    // change "unix" to "tcp" to use over a network.
    os.Remove(me) // only needed for "unix"
    l, e := net.Listen("unix", me);
    if e != nil {
        log.Fatal("listen error: ", e);
    }
    ls.l = l

    // please don't change any of the following code,
    // or do anything to subvert it.

    // create a thread to accept RPC connections from clients.
    go func() {
        for ls.dead == false {
            conn, err := ls.l.Accept()
            if err == nil && ls.dead == false {
                if ls.dying {
                    // process the request but force discard of reply.

                    // without this the connection is never closed,
                    // b/c ServeConn() is waiting for more requests.
                    // test_test.go depends on this two seconds.
                    go func() {
                        time.Sleep(2 * time.Second)
                        conn.Close()
                    }()
                    ls.l.Close()

                    // this object has the type ServeConn expects,
                    // but discards writes (i.e. discards the RPC reply).
                    deaf_conn := DeafConn{c : conn}

                    rpcs.ServeConn(deaf_conn)

                    ls.dead = true
                } else {
                    go rpcs.ServeConn(conn)
                }
            } else if err == nil {
                conn.Close()
            }
            if err != nil && ls.dead == false {
                log.Printf("LockServer(%v) accept: %v\n", me, err.Error())
                ls.kill()
            }
        }
    }()

    return ls
}
