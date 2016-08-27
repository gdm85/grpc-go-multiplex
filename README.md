# grpc-go-multiplex

This is a working example/proof-of-concept that shows how to use [soheilhy's excellent cmux package](http://github.com/soheilhy/cmux) (cockroachdb's fork is used here).

This example covers a very specific use-case, as per [release blog post](https://medium.com/@gdm85/http-load-balancing-on-grpc-services-e3d702db05d7); for more use-cases give a look to the examples in [cmux' GoDocs](https://godoc.org/github.com/soheilhy/cmux).

For upstream discussion, see https://github.com/grpc/grpc-go/issues/549.

# License

[GNU GPL version 2](./LICENSE)
The greeter client side example is covered by a different BSD-style license (see [greeter_client.go](./greeter_client/greeter_client.go)).

# How to build

This project uses an automatically-provisioned GOPATH. Example init/building commands on a Linux system:

```
git submodule update --init --recursive
make
```

# How to use

Start the server-side in a terminal:
```
$ bin/greeter_multiplex_server
2016/08/27 16:20:21 listening and serving (multiplexed) on :50051
```

Make a gRPC test in another terminal with:
```
$ bin/greeter_client
2016/08/27 16:21:47 Greeting: Hello world
```

Test a status call (for example, as load-balancers would do) by fetching `/status`:
```
$ curl -v http://localhost:50051/status
* Hostname was NOT found in DNS cache
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 50051 (#0)
> GET /status HTTP/1.1
> User-Agent: curl/7.50.1
> Host: localhost:50051
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Sat, 27 Aug 2016 14:24:02 GMT
< Content-Length: 2
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host localhost left intact
OK
```

And browse any other URL (via browser) to see normal HTML output:
```
$ curl http://localhost:50051/xxx
Welcome to the home page!
```
