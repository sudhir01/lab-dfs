package lockservice
import "crypto/rand"
import "math/big"

//
// RPC definitions for a simple lock service.
//
// You will need to modify this file.
//

//
// Lock(lockname) returns OK=true if the lock is not held.
// If it is held, it returns OK=false immediately.
// 
type LockArgs struct {
  // Go's net/rpc requires that these field
  // names start with upper case letters!
  Lockname      string  // lock name
  RequestId     int64   // unique id
  RequestSource string  // who is making this request
}

type LockReply struct {
  OK bool
}

//
// Unlock(lockname) returns OK=true if the lock was held.
// It returns OK=false if the lock was not held.
//
type UnlockArgs struct {
  Lockname      string
  RequestId     int64   //unique identifier
  RequestSource string  // who is making this request
}

type UnlockReply struct {
  OK bool
}

func randomId() int64 {
    max := big.NewInt(int64(1) << 62)
    bigx, _ := rand.Int(rand.Reader, max)
    x := bigx.Int64()
    return x
}
