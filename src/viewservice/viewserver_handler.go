package viewservice

type ViewServerHandler struct {
}

// 
// server Get() RPC handler.
//
func(handler *ViewServerHandler) Get(args *GetArgs, reply *GetReply) error {

  // Your code here.

  return nil
}

//TODO - replace with the Ping method from server.go
func(handler *ViewServerHandler) Ping(args *PingArgs, reply *PingReply) error {
	 return nil
}
