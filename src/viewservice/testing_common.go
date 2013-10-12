package viewservice

import "os"
import "strconv"

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

func InitServer(serverHostPort string) *ViewServer {
    server := StartServer(serverHostPort)
    return server
}

func InitClient(clientPortName string, serverHostPort string) *Clerk {
    client := MakeClerk(Port(clientPortName), serverHostPort)
    return client
}

