package viewservice

import "errors"
import "net/rpc"

/**
ViewServerFactory initializes a ViewServer. It will be useful to separate the logic that handles RPC requests in the ViewServer from the ViewServer itself. This will allow testing the dispatch logic independent of the server initialization code.
 */
func NewViewServer(hostPort string, rpcServer *rpc.Server, tracker *ViewTracker, handler ServerHandler) (*ViewServer, error) {
	 err := validateInput(hostPort, rpcServer, handler)
	 if err != nil {
		  return nil, err
	 }

    vs := new(ViewServer)

    vs.me          = hostPort
    vs.rpcServer   = rpcServer
	 vs.handler     = handler
	 vs.tracker	    = tracker
    return vs, nil
}


func validateInput(hostPort string, rpcServer *rpc.Server, handler ServerHandler) error {
	 if hostPort == "" {
		  return errors.New("hostPort cannot be empty")
	 }
	 if rpcServer == nil {
		  return errors.New("rpc server cannot be nil")
	 }
	 if handler == nil {
		  return errors.New("server handler cannot be nil")
	 }
	 return nil
}
