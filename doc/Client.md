Client
- Call
  // async create a Call invocation
- Go


### Call: All info needed for one RPC call.
```go
// Call represents an active RPC.
type Call struct {
	Seq           uint64
	ServiceMethod string      // format "<service>.<method>"
	Args          interface{} // arguments to the function
	Reply         interface{} // reply from the function
	Error         error       // if error occurs, it will be set
	Done          chan *Call  // Strobes when call is complete.
}

func (call *Call) done() {
	call.Done <- call
}
```
Done provide the interface to notify the caller that this Call has finished. And the reply has been encoded back to the **Reply**



### Client

```go
// Client represents an RPC Client.
// There may be multiple outstanding Calls associated
// with a single Client, and a Client may be used by
// multiple goroutines simultaneously.
type Client struct {
	cc       codec.Codec
	opt      *Option
	sending  sync.Mutex // protect following
	header   codec.Header
	mu       sync.Mutex // protect following/
	seq      uint64
	pending  map[uint64]*Call
	closing  bool // user has called Close
	shutdown bool // server has told us to stop
}
```
**cc**: Critical component, as the conn is a member of it, this is like an IO, supporting encryption for the RPC
**opt**: 1st step, json for encoding the Connection Option, LIKE how encoding algo and magic number, decide whether client can connected to the server, the Auth wiil also be part of it in future. 
**sending**: to protect the pending Call Map


#### Make a Call
```go
func (client *Client) Call(serviceMethod string, args, reply interface{}) error 
{
	-> Go
}
func (client *Client) Go(serviceMethod string, args, reply interface{}, done chan *Call) *Call 
{
	Construct Call
	send(Call)
}

send(call){
	client.registerCall(call)
	    // give get new sequence number
		// save to pending Call Map
	client.cc.Write(client.header, call.Args)
	    // actually write to the connections
    call := client.removeCall(seq)
	call.done() 
        // Close call;
}
```
Go invokes the function asynchronously.
It returns the Call structure representing the invocation.



#### Client Initialization
```go
func newClientCodec(cc codec.Codec, opt *Option) *Client {
    client := &Client{
        seq:     1, // seq starts with 1, 0 means invalid call
        cc:      cc,
        opt:     opt,
        pending: make(map[uint64]*Call),
    }
	go client.receive()
    return client
}
```



