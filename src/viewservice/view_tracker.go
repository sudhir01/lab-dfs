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

func elapsedDeadPings(lastPing time.Time) int64 {
    now   := time.Now()
    delta := now.Sub(lastPing)
    return int64(delta/PingInterval)
    //milli := (delta.Nanoseconds()/1000)
    //pings := milli / (PingInterval.Nanoseconds()/1000)
    //return pings
}

/* tick() is called once per PingInterval.
- thread safe
- It should notice
if servers have died or recovered, and change the view
accordingly.
Periodic tasks on ping:
1. Mark server dead if max DeadPings have passed for PingIntervals
2. Update view for either (only if primary has not drifted):
i.  dead server or                      -> TODO
ii. idle server when there is no backup -> TODO
*/
func (this  *ViewTracker) tick() {
  // Your code here.
}
