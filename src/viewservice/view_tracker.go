package viewservice

import (
	 "sync"
	 "time"
)

type ViewTracker struct {
  mu          sync.Mutex
  pingTimes   map[string] time.Time
  currentView View
}

func NewViewTracker() *ViewTracker {
	 this := new(ViewTracker)
	 this.pingTimes   = map[string] time.Time{}
	 this.currentView = View{INITIAL_VIEW, NO_SERVER, NO_SERVER, NO_VIEW, NO_VIEW}
	 return this
}

func (this *ViewTracker) Ping(args *PingArgs, reply *PingReply) error {
	 this.mu.Lock()
	 defer this.mu.Unlock()

	 viewnum                := args.Viewnum
	 server                 := args.Me
	 this.pingTimes[server]  = time.Now()

	 switch server {
	 case this.currentView.Primary:
		  this.currentView.PrimaryView = viewnum
	 case this.currentView.Backup:
		  this.currentView.BackupView  = viewnum
	 }
	 reply.View = this.currentView
	 return nil
}

func (this *ViewTracker) PingTable() *map[string] time.Time {
    return &this.pingTimes
}

func (this *ViewTracker) View() *View {
    return &this.currentView
}
