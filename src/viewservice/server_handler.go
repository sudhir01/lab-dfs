package viewservice

type ServerHandler interface {
	 Get(args *GetArgs, reply *GetReply) error
	 Ping(args *PingArgs, reply *PingReply) error
}
