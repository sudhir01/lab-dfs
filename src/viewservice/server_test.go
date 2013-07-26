package viewservice

import "testing"
// import "runtime"
// import "time"
// import "fmt"
// import "os"

func Test_init_view_server(t *testing.T) {
    hostPort := Port("v")

    noserver, err := NewViewServer("")
    if noserver != nil {
        t.Fatalf("Server was created when no hostname was provided\n")
    }

    if err == nil {
        t.Fatalf("Error message not returned for invalid server initialization parameters\n")
    }

    server, err := NewViewServer(hostPort)
    if server == nil {
        t.Fatalf("Could not initialize view server. Server reference is nil\n")
    }

    if err != nil {
        t.Fatalf("Could not initialize view server. Got error %s\n", err.Error())
    }

    if server.IsDead() {
        t.Fatalf("New server marked dead\n", err.Error())
    }

    if server.Name() != hostPort {
        t.Fatalf("Server was not initialzied with host name: [%s]\n", hostPort)
    }

    if server.ListenerAddress() != hostPort {
        t.Fatalf("Server is not listening on the host port [%s]\n", hostPort)
    }

    expectedView := &View{INITIAL_VIEW, NO_SERVER, NO_SERVER, NO_VIEW, NO_VIEW}
    actualView   := server.View()
    if actualView != expectedView {
        t.Fatalf("Server was not initialized with expectedView [+%v], got view [+%v]\n", expectedView, actualView)
    }

    pingTable := server.PingTable()
    if pingTable == nil || *pingTable == nil || len(*pingTable) != 0 {
        t.Fatalf("Server's ping table is not empty, pingTable [+%v]\n", pingTable)
    }
}

func Test_first_view(t *testing.T) {
    t.FailNow()
}

func Test_primary_is_primary_from_last_view(t *testing.T) {
    t.FailNow()
}

func Test_primary_is_backup_from_last_view(t *testing.T) {
    t.FailNow()
}

func Test_any_primary_on_view_server_init(t *testing.T) {
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

