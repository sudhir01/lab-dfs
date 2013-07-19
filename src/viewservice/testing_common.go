package viewservice

import "os"
import "strconv"

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

