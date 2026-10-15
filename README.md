# go-grpc-interceptor

gRPC server interceptors for [grpc-go](https://github.com/grpc/grpc-go).

# Interceptors

| interceptor | Description |
| -----------|--------|
| [multiinterceptor](https://github.com/eltorocorp/go-grpc-request-id-interceptor/tree/master/multiinterceptor) | Chain multiple `UnaryServerInterceptor`s or `StreamServerInterceptor`s |
| [panichandler](https://github.com/eltorocorp/go-grpc-request-id-interceptor/tree/master/panichandler) | Protect a process from aborting by panic and return Internal error as status code |
| [xrequestid](https://github.com/eltorocorp/go-grpc-request-id-interceptor/tree/master/xrequestid) | Generate X-Request-Id to each request |
| [acceptlang](https://github.com/eltorocorp/go-grpc-request-id-interceptor/tree/master/acceptlang) | Parses Accept-Language from metadata |
| [instrument](https://github.com/eltorocorp/go-grpc-request-id-interceptor/tree/master/instrument) | Instrumentation hook |

# Committers

* Masahiro Sano(@kazegusuri)

# Contribution

Please read the CLA below carefully before submitting your contribution.

<https://www.mercari.com/cla/>

# License

Copyright 2016 Mercari, Inc.

Licensed under the MIT License.
