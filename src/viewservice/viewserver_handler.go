package viewservice

import (
	 "sync"
	 "time"
)

type ViewServerHandler struct {
  mu          sync.Mutex
  pingTimes   map[string] time.Time
  currentView View
}

func NewViewServerHandler() *ViewServerHandler {
	 handler := new(ViewServerHandler)
	 handler.pingTimes   = map[string] time.Time{}
	 handler.currentView = View{INITIAL_VIEW, NO_SERVER, NO_SERVER, NO_VIEW, NO_VIEW}
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
	 handler.mu.Lock()
	 defer handler.mu.Unlock()

	 viewnum                   := args.Viewnum
	 server                    := args.Me
	 handler.pingTimes[server] = time.Now()

	 switch server {
	 case handler.currentView.Primary:
		  handler.currentView.PrimaryView = viewnum
	 case handler.currentView.Backup:
		  handler.currentView.BackupView  = viewnum
	 }

	 reply.View = handler.currentView
	 return nil
}

//FIXME - this method cannot be exported - remove from handler
func (handler *ViewServerHandler) PingTable() *map[string] time.Time {
    return &handler.pingTimes
}

//FIXME - this method cannot be exported - remove from handler
func (handler *ViewServerHandler) View() *View {
    return &handler.currentView
}
