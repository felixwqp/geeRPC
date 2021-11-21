
For a RPC call:
```go
err = client.Call("Service.Method", args, &reply)
```

Client Request with:
1. Service Name
2. Method Name
3. Arguments

Service Response: body & header
body: args, return val
header: other info, 

interface for serialization

