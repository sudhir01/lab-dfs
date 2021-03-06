package lockservice

import "net/rpc"
import "log"

//
// the lockservice Clerk lives in the client
// and maintains a little state.
//
type Clerk struct {
    servers [2]string // primary port, backup port
    // TODO Your definitions here.
}


func MakeClerk(primary string, backup string) *Clerk {
    ck := new(Clerk)
    ck.servers[0] = primary
    ck.servers[1] = backup
    // TODO Your initialization code here.
    return ck
}

/*
* call() sends an RPC to the rpcname handler on server srv
* with arguments args, waits for the reply, and leaves the
* reply in reply. the reply argument should be the address
* of a reply structure.
*
* call() returns true if the server responded, and false
* if call() was not able to contact the server. in particular,
* reply's contents are valid if and only if call() returned true.
*
* you should assume that call() will time out and return an
* error after a while if it doesn't get a reply from the server.
*
* please use call() to send all RPCs, in client.go and server.go.
* please don't change this function.
*/
func call(server string, rpcname string,
          args interface{}, reply interface{}) bool {
    connection, errx := rpc.Dial("unix", server)
    if errx != nil {
        return false
    }
    defer connection.Close()

    err := connection.Call(rpcname, args, reply)
    if err == nil {
        return true
    }
    return false
}

/*
* calls call() for the first server in the list of servers. If it is unable to
* reach the server, then it tries to call the next server in the array.
*/
func callWithFallback(servers []string,
                      rpcname string,
                      args    interface{},
                      reply   interface{}) bool {
    called      := false
    serverCount := len(servers)
    tries       := 0

    log.Printf("[debug] Client::callWithFallback arguments serverCount: %v \t rpcname: %v \t args: %v", serverCount, rpcname, args)
    for called == false && tries < serverCount {
        server := servers[tries]
        tries += 1
        log.Printf("[debug] Client::callWithFallback server(%v) attempt(%v) args(%v) calling..", server, tries, args)
        called = call(server, rpcname, args, reply)
        log.Printf("[debug] Client::callWithFallback server(%v) attempt(%v) response(%v)", server, tries, called)
    }

    return called
}


//
// ask the lock service for a lock.
// returns true if the lock service
// granted the lock, false otherwise.
//
// TODO you will have to modify this function.
//
func (ck *Clerk) Lock(lockname string) bool {
    // prepare the arguments.
    var requestId int64 = randomId()
    args := &LockArgs{lockname, requestId, "client"}
    args.Lockname = lockname
    var reply LockReply

    // send an RPC request, wait for the reply.
    //TODO - handle the case where we are unable to contact the server
    ok := callWithFallback(ck.servers[0:], "LockServer.Lock", args, &reply)
    if ok == false {
        return false
    }

    return reply.OK
}


//
// ask the lock service to unlock a lock.
// returns true if the lock was previously held,
// false otherwise.
// TODO - implement this function

func (ck *Clerk) Unlock(lockname string) bool {
    // prepare the arguments.
    var requestId int64 = randomId()
    args := &UnlockArgs{lockname, requestId, "client"}
    var reply UnlockReply

    //ask the lock service to unlock
    ok := callWithFallback(ck.servers[0:], "LockServer.Unlock", args, &reply)
    if ok == false {
        //TODO - handle the case where we are unable to contact the server
        return false
    }
    return reply.OK
}
