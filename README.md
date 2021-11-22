# RPC framework

RPC framework 
1. supporting concurrent C/S, 
2. with HTTP transmission, 
3. encoding message by gob, 
4. supporting Timeout with contextTimeout/channel/net.TimeoutDial
5. supporting load balancing and heartbeating
6. registry for server discovery. 

More details: 
- [Client](doc/Client.md)
- [Http Support](doc/httpSupport.md)
- [MessageEncoding/Decoding in Conn](doc/message_serialization_deseriailization.md)
- [Registry](doc/Registry.md)
- [Timeout Support](doc/TimeoutSupport.md)
- [LoadBalance Support](doc/LoadBalance.md)
- [Service Lookup via Reflection](doc/ServiceRegister.md)


Start the demo
```shell
make run
```


reference:  
- [GeeRPC](https://geektutu.com/post/geerpc.html)
- [golang net/rpc](https://pkg.go.dev/net/rpc)
- [rpcX](https://github.com/smallnest/rpcx)