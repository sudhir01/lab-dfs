package viewservice

import (
	 "testing"
	 "net/rpc"
	 "reflect"
)

func Test_server_accepts_connection_for_get_and_ping(t *testing.T) {
	 hostPort  := Port("v")
	 rpcServer := rpc.NewServer()
	 handler   := new(TestServerHandler)
	 timer	  := new(DefaultTimer)
	 tracker	  := NewViewTracker(timer)
	 server, _ := NewViewServer(hostPort, rpcServer, tracker, handler)

	 server.Start()
	 clnt, errx := rpc.Dial("unix", hostPort)
	 if errx != nil {
		  t.Fatalf("Unable to establish a connection wit the server %s, got error %s\n", hostPort, errx.Error())
	 }
	 defer clnt.Close()
	 getReply := new(GetReply)
	 getArgs  := new(GetArgs)

	 clnt.Call("ViewServer.Get", getArgs, getReply)
	 actualView   := getReply.View
	 expectedView := View{1, "primary", "backup", 1, 1}
    if reflect.DeepEqual(actualView, expectedView) == false {
		  t.Fatalf("Server Get reply (+%v) does not match expected reply (+%v)", actualView, expectedView)
	 }
}

func Test_init_view_server(t *testing.T) {
    hostPort  := Port("v")
    rpcServer := rpc.NewServer()
	 timer	  := new(DefaultTimer)
	 tracker	  := NewViewTracker(timer)
	 handler   := NewViewServerHandler(tracker)

    noServer, err := NewViewServer("", rpcServer, tracker, handler)
    if noServer != nil {
        t.Fatalf("Server was created when no hostname was provided\n")
    }

    if err == nil {
        t.Fatalf("Error message not returned for invalid server initialization parameters\n")
    }

    noServer, err = NewViewServer("test-port", nil, tracker, handler)
    if noServer != nil {
        t.Fatalf("ViewServer was created when the RPC server was nil\n")
    }

    if err == nil {
        t.Fatalf("NewViewServer did not return an error code when the RPC server was nil\n")
    }

    server, err := NewViewServer(hostPort, rpcServer, tracker, handler)
    if server == nil {
        t.Fatalf("Could not initialize view server. Server reference is nil\n")
    }

    if err != nil {
        t.Fatalf("Could not initialize view server. Got error %s\n", err.Error())
    }

    if server.IsListening() == true {
        t.Fatalf("New server is listening before starting the server\n")
    }

    server.Start()

    if server.IsListening() != true {
        t.Fatalf("New server is not listening after starting\n")
    }

    if server.IsDead() {
        t.Fatalf("New server marked dead\n", err.Error())
    }

    if server.Name() != hostPort {
        t.Fatalf("Server was not initialzied with host name: [%s]\n", hostPort)
    }

    listenerAddr := server.ListenerAddress()
    if listenerAddr != hostPort {
        t.Fatalf("Server is not listening on the host port expected[%s], actual [%s]\n", hostPort, listenerAddr)
    }

    //TODO - add a test to ensure that the connection acceptor is running
    //TODO - add a test to ensure that the ticker is running
}
