package viewservice

type ViewServerHandler struct {
	 tracker *ViewTracker
}

func NewViewServerHandler(tracker *ViewTracker) *ViewServerHandler {
	 handler := new(ViewServerHandler)
	 handler.tracker = tracker
	 return handler
}

// 
// server Get() RPC handler.
//
func(handler *ViewServerHandler) Get(args *GetArgs, reply *GetReply) error {

  // Your code here.

  return nil
}

func (handler *ViewServerHandler) Ping(args *PingArgs, reply *PingReply) error {
	 handler.tracker.Ping(args, reply)
	 return nil
}
