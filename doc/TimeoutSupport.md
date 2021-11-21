

Client Timeout:
- Connect to Server address: session connection
  -  Dial to a certain session
- Send Request: fail to encode and write to Conn
- Receive from the server
- Read message: fail to decode

Server Timeout:
- read request message
- send response message
- call the service(reflect)

In This project:
1. client create connection
   - new.Dial -> net.DialTimeout
   - NewClient, options communication, finishes,
2. Client.Call(), includng read/write
3. serve handle request.
   - 




```go
func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	// when the service call actually happens!
	// if the call never return, this goroutine will die here
	err := req.svc.call(req.mtype, req.argv, req.replyv)
	if err != nil {
		req.h.Error = err.Error()
		server.sendResponse(cc, req.h, invalidRequest, sending)
		return
	}
	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}
```

add two channel **called** **sent**, to make sure sendResponse(), can be called if and only if once.