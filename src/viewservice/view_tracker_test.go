package viewservice

import (
	 "testing"
	 "time"
)

func Test_tracker_initialization(t *testing.T) {
	 timer        := &MockTimer{0, 0}
	 tracker	     := NewViewTracker(timer)
	 expectedView := &View{INITIAL_VIEW, NO_SERVER, NO_SERVER, INITIAL_VIEW, INITIAL_VIEW}
	 checkView(tracker, expectedView, t)

	 pingTable := tracker.PingTable()
	 if pingTable == nil || pingTable == nil || len(pingTable) != 0 {
		  t.Fatalf("Server's ping table is not empty, pingTable [+%v]\n", pingTable)
	 }
}

func Test_server_becomes_primary_on_first_ping_after_initialization(t *testing.T) {
	 server1, timer, tracker, ping, reply := ViewTrackerVariables()

	 tracker.Ping(ping, reply)

	 //check view
	 expectedView := &View{1, server1, NO_SERVER, INITIAL_VIEW, INITIAL_VIEW}
	 checkView(tracker, expectedView, t)

	 //check ping table
	 time1 := timer.Now()
	 expectedTable := map[string] time.Time { server1: time1}
	 checkTable(tracker, expectedTable, t)
}

func Test_view_does_not_change_till_primary_acks(t *testing.T) {
	 //server1, timer, tracker, ping, reply := ViewTrackerVariables()
}

func Test_primary_is_primary_from_last_view(t *testing.T) {
	 //server1 := "server-1-primary"
	 //server2 := "server-2-backup"
	 //server3 := "server-3-idle"
    t.FailNow()
}

func Test_primary_is_backup_from_last_view(t *testing.T) {
    t.FailNow()
}

func Test_only_one_primary(t *testing.T) {
    t.FailNow()
}

func Test_idle_server_becomes_backup(t *testing.T) {
    t.FailNow()
}

func Test_primary_does_not_become_backup(t *testing.T) {
    t.FailNow()
}

func Test_replies_to_ping_with_view(t *testing.T) {
    t.FailNow()
}

func Test_change_view_when_primary_is_dead_with_backup(t *testing.T) {
    t.FailNow()
}

func Test_change_view_when_primary_is_dead_without_backup(t *testing.T) {
    t.FailNow()
}

func Test_change_view_when_backup_is_dead(t *testing.T) {
    t.FailNow()
}

//mark server dead
//

func Test_mark_primary_dead(t *testing.T) {
    t.FailNow()
}

func Test_mark_backup_dead(t *testing.T) {
    t.FailNow()
}

//keep view in sync with what primary sees

func Test_view_does_not_change_till_primary_ack_idle_server(t *testing.T) {
    t.FailNow()
}

func Test_view_does_not_change_till_primary_ack_backup_dead(t *testing.T) {
    t.FailNow()
}

func Test_view_does_not_change_till_primary_ack_primary_dead(t *testing.T) {
    t.FailNow()
}

func Test_view_does_not_change_till_primary_ack_primary_and_backup_dead(t *testing.T) {
    t.FailNow()
}
