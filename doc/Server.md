
```go
// request stores all information of a call
type request struct {
	h            *codec.Header // header of request
	argv, replyv reflect.Value // argv and replyv of request
	mtype        *methodType
	svc          *service
}
```

```go
	req.argv = req.mtype.newArgv()
	req.replyv = req.mtype.newReplyv()

	// make sure that argvi is a pointer, ReadBody need a pointer as parameter
	argvi := req.argv.Interface()
	if req.argv.Type().Kind() != reflect.Ptr {
		argvi = req.argv.Addr().Interface()
	}
	if err = cc.ReadBody(argvi); err != nil {
		log.Println("rpc server: read body err:", err)
		return req, err
	}
```

readRequest 方法中最重要的部分，即通过 newArgv() 和 newReplyv() 两个方法创建出两个入参实例，然后通过 cc.ReadBody() 将请求报文反序列化为第一个入参 argv，在这里同样需要注意 argv 可能是值类型，也可能是指针类型，所以处理方式有点差异。



server add the service to ServiceMap during the registration.