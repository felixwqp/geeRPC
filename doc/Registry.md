![](Registry.png)

Client and Serser don't need to know each other, they only needs to know the registry
1. On stratup, 
- Server register itself on the Registry; Registry mark its service as available; 
```go
func startServer(registryAddr string, wg *sync.WaitGroup) {
	var foo Foo
	l, _ := net.Listen("tcp", ":0")
	server := geerpc.NewServer()
	_ = server.Register(&foo)
	registry.Heartbeat(registryAddr, "tcp@"+l.Addr().String(), 0)
	wg.Done()
	server.Accept(l)
}
```
goroutine + channel for periodic sending
```go
// Heartbeat send a heartbeat message every once in a while
// it's a helper function for a server to register or send heartbeat
func Heartbeat(registry, addr string, duration time.Duration) {
	if duration == 0 {
		// make sure there is enough time to send heart beat
		// before it's removed from registry
		duration = defaultTimeout - time.Duration(1)*time.Minute
	}
	var err error
	err = sendHeartbeat(registry, addr)
	go func() {
		t := time.NewTicker(duration)
		for err == nil {
			<-t.C
			err = sendHeartbeat(registry, addr)
		}
	}()
}
```
- Server sends Heartbeats to Registry to update the availablity
2. Client:
- Ask Registry to current available service, registry will return the available server address list
- With server address list, client make a RPC call based on the load balance algorithm


```go
type RegistryDiscovery struct {
	*MultiServersDiscovery
	registry   string
	timeout    time.Duration 	// refresh interval
	lastUpdate time.Time
}
```
registry is in part of the Registry Discovery




Registry is serving HTTP request from server/clients on a certain address;
- for GET from the client for available server
- for POST from server for heartbeat updating
```go
// Runs at /_geerpc_/registry
func (r *GeeRegistry) ServeHTTP(w http.ResponseWriter, req *http.Request) {
switch req.Method {
case "GET":
// keep it simple, server is in req.Header
w.Header().Set("X-Geerpc-Servers", strings.Join(r.aliveServers(), ","))
case "POST":
// keep it simple, server is in req.Header
addr := req.Header.Get("X-Geerpc-Server")
if addr == "" {
w.WriteHeader(http.StatusInternalServerError)
return
}
r.putServer(addr)
default:
w.WriteHeader(http.StatusMethodNotAllowed)
}
}

// HandleHTTP registers an HTTP handler for GeeRegistry messages on registryPath
func (r *GeeRegistry) HandleHTTP(registryPath string) {
http.Handle(registryPath, r)
log.Println("rpc registry path:", registryPath)
}

```





XClient will inquery the discovery to get the latest available server, 
```go
// Call invokes the named function, waits for it to complete,
// and returns its error status.
// xc will choose a proper server.
func (xc *XClient) Call(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	rpcAddr, err := xc.d.Get(xc.mode)
	if err != nil {
		return err
	}
	return xc.call(rpcAddr, ctx, serviceMethod, args, reply)
}
```