package viewservice

import (
	 "sync"
	 "time"
)

type ViewTracker struct {
  mu          sync.Mutex
  pingTimes   map[string] time.Time
  currentView View
  timer       Timer
}

func NewViewTracker(timer Timer) *ViewTracker {
	 this := new(ViewTracker)
	 this.pingTimes   = map[string] time.Time{}
	 this.currentView = View{INITIAL_VIEW, NO_SERVER, NO_SERVER, INITIAL_VIEW, INITIAL_VIEW}
	 this.timer       = timer
	 return this
}

func (this *ViewTracker) Ping(args *PingArgs, reply *PingReply) error {
	 this.mu.Lock()
	 defer this.mu.Unlock()

	 viewnum                := args.Viewnum
	 server                 := args.Me
	 this.pingTimes[server] = this.timer.Now()

	 //the view server booted, first ping
	 if this.currentView.Viewnum == 0  && viewnum == 0 {
		  this.currentView.Viewnum = 1
		  this.currentView.Primary = server
	// view server rebooted, old server still alive
	 } else if this.currentView.Viewnum != 0 && viewnum == 0 {
	 } else if this.currentView.Viewnum == 0 && viewnum != 0 {
	 } else if this.currentView.Viewnum != 0 && viewnum != 0 {
	 }
		  //the server crashed or booted
		  /*
		  if viewnum == 0 {
				if server == this.currentView.Primary {
				} else {
				}
		  } else if viewnum == 0 && server == this.currentView.Backup {

		  } else if viewnum > 0 && server == this.currentView.Primary {
		  } else if viewnum > 0 && server == this.currentView.Backup {
		  } else if viewnum  {}
		  */
	 reply.View = this.currentView
	 return nil
}

func (this *ViewTracker) PingTable() map[string] time.Time {
    return this.pingTimes
}

func (this *ViewTracker) View() *View {
    return &this.currentView
}
/* tick() is called once per PingInterval.
- thread safe
- It should notice if servers have died or recovered, and change the view accordingly.

Periodic tasks on ping:
1. Mark server dead if max DeadPings have passed for PingIntervals
2. Update view for either (only if primary has not drifted):
i.  dead server or                      -> TODO
ii. idle server when there is no backup -> TODO
*/
func (this  *ViewTracker) Tick() {
  // Your code here.
}

//private methods

func (this *ViewTracker) ElapsedDeadPings(lastPing time.Time) int64 {
    now   := this.timer.Now()
    delta := now.Sub(lastPing)
    return int64(delta/PingInterval)
    //milli := (delta.Nanoseconds()/1000)
    //pings := milli / (PingInterval.Nanoseconds()/1000)
    //return pings
}
