package viewservice

import (
	 "os"
	 "strconv"
	 "testing"
	 "time"
	 "reflect"
)

type TestServerHandler struct {}
func (TestServerHandler) Get(args *GetArgs, reply *GetReply) error {
	 view := View{1, "primary", "backup", 1, 1}
	 reply.View = view
	 return nil
}
func (TestServerHandler) Ping(args *PingArgs, reply *PingReply) error {
	 return nil
}

func Port(suffix string) string {
  s := "/var/tmp/824-"
  s += strconv.Itoa(os.Getuid()) + "/"
  os.Mkdir(s, 0777)
  s += "viewserver-"
  s += strconv.Itoa(os.Getpid()) + "-"
  s += suffix
  return s
}

func checkView(tracker *ViewTracker, expectedView *View, t *testing.T) {
	 actualView   := tracker.View()
	 if reflect.DeepEqual(actualView, expectedView) == false {
		  t.Fatalf("Tracker expected view [+%v], got view [+%v]\n", expectedView, actualView)
	 }
}

func checkTable(tracker *ViewTracker, expectedTable map[string] time.Time, t *testing.T) {
	 actualTable := tracker.PingTable()

	 if reflect.DeepEqual(actualTable, expectedTable) == false {
		  t.Fatalf("Tracker expected ping table [+%v], got ping table [+%v]\n", expectedTable, actualTable)
	 }
}

type MockTimer struct {
	 sec  int64
	 nsec int64
}

func (this *MockTimer) Decrement(sec int64, nsec int64) {
	 this.sec -= sec
	 this.nsec -= nsec
}

func (this *MockTimer) Increment(sec int64, nsec int64) {
	 this.sec -= sec
	 this.nsec -= nsec
}

func (this *MockTimer) Now() time.Time {
	 return time.Unix(this.sec, this.nsec)
}

func ViewTrackerVariables() (string, Timer, *ViewTracker, *PingArgs, *PingReply, *View) {
	 primary	     := "server-1-primary"
	 timer		  := &MockTimer{0, 0}
	 tracker	     := NewViewTracker(timer)
	 ping         := &PingArgs{primary, INITIAL_VIEW}
	 reply        := new(PingReply)
	 expectedView := &View{1, primary, NO_SERVER, INITIAL_VIEW, INITIAL_VIEW}

	 return primary, timer, tracker, ping, reply, expectedView
}
