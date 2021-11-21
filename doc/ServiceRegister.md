
rpc needs to have format like
```go
func (t *T) MethodName(argType T1, replyType *T2) error
```
to be eligible for a rpc call, specifically:
```md
the method’s type is exported. (Capitalized)
the method is exported. (Capitalized)
the method has two arguments, both exported (or builtin) types. 
the method’s second argument is a pointer. // points to 
the method has return type error. 
```

```go
type methodType struct {
	method 		reflect.Method
	ArgType 	reflect.Type
	ReplyType	reflect.Type
	numCalls	uint64
}
func (m *methodType) newArgv() reflect.Value {
func (m *methodType) newReplyv() reflect.Value {
```
the funcs provide the interface to initialize the Args/Reply with corresponding type.
these type can be take by the interface methods
```go
	argType, replyType := mType.In(1), mType.In(2)
		if !isExportedOrBuiltinType(argType) || !isExportedOrBuiltinType(replyType) {
			continue
		}
		s.method[method.Name] = &methodType{
			method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}
```


```go
s := newService(rcvr)
	
type service struct {
    name   string
    typ    reflect.Type
    rcvr   reflect.Value
    method map[string]*methodType
}

```


### interface for users
write to the request, create the request with the methodType; with args/reply/method
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
    log.Println("rpc server: read argv err:", err)
    return req, err
}
```

	err := req.svc.call(req.mtype, req.argv, req.replyv)


reflection is a key component of this service register feature.
reflect has two key type: **reflect.Type** and  **reflect.Value**
### reflect.Type && reflect.Value
-> just a Go Type, this is an interface with many methods
, implmentation is type descriptor
1. discriminate among types
2. inspect the components
```go
t := reflect.TypeOf(3)  // a reflect.Type
fmt.Println(t.String()) // "int"
fmt.Println(t)          // "int"
```
**3 -> interface{}**: 
an assignment from a concrete value to an interface type performs an **implicit interface** conversion, which creates an interface value consisting of two components:
- its dynamic type is the operand’s type (int)
- its dynamic value is the operand’s value (3).


Typeof and Valueof expose these dynamic Type:
- Typeof return reflect.Type. an interface value’s dynamic type, it always returns a concrete type
- Valueof return reflect.Value,  containing the interface’s dynamic value.


The inverse operation to reflect.ValueOf is the reflect.Value.Interface method. It returns an interface{} holding the same concrete value as the reflect.Value:

```go
v := reflect.ValueOf(3) // a reflect.Value
x := v.Interface()      // an interface{}
i := x.(int)            // an int
fmt.Printf("%d\n", i)   // "3"
```


### reflect.Value() v.s. interface{}
A reflect.Value and an interface{} can both hold arbitrary values. 
- an empty interface hides the representation and intrinsic operations of the value it holds and exposes none of its methods, so unless we know its dynamic type and use a type assertion to peer inside it (as we did above), there is little we can do to the value within. 
- In contrast, a Value has many methods for inspecting its contents, regardless of its type.
```go
s.rcvr.Method(0).String()
result = {string} "<func(geerpc.Args, *int) error Value>"

s.rcvr.Interface()
result = {interface {} | github.com/felixwqp/geerpc.Foo}
Random1 = {int} 0
Random2 = {int} 0
```



