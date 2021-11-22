


decouple the multiple service instance and the server; Discovery manages the server address for **one** certain service;

SelectMode

```go
type Discovery interface {
	Refresh() error // refresh from remote registry
	Update(servers []string) error
	Get(mode SelectMode) (string, error)
	GetAll() ([]string, error)
}
```


MultiServersDiscovery is a discovery for multi servers without a registry center user provides the server addresses explicitly instead
which implements the Discovery, 
```go
type MultiServersDiscovery struct {
	r       *rand.Rand   // generate random number
	mu      sync.RWMutex // protect following
	servers []string
	index   int // record the selected position for robin algorithm
}
```
servers is temporarily represents by string of address like "tcp@addr"
Discovery will get the address of a certain server by the select mode; 


### Xclient
xclient: interface with service discovery;
One XClient is used only for one service. You should create multiple XClient for multiple services.

each
```go
type XClient struct {
d       Discovery
mode    SelectMode
opt     *geerpc.Option
mu      sync.Mutex // protect following
clients map[string]*geerpc.Client
}
Call(ctx, serviceMethod, args, reply) -> error
Broadcast(ctx, serviceMethod, args, relpy) -> error
    // after geting the server address
    xc.call(rpcAddr, ctx, serviceMethod, args, reply)
        xc.dial(rpcAddr) -> Client
            check whether the current client on rpcAddr is available, if so return, o.w. create one.
        Client.Call(ctx, serviceMethod, args, reply)
```

every rpc address will have one client to make call. 
every Call/Broadcast will select the address from the discovery service with configured SelectMode; 





