# Serve It All

A simple HTTP/1.1, HTTP/2 and grpc server to be used as a backend for Istio.

## Usage

```shell
PORT=8080 ./main
```

Via HTTP/1.1
```shell
curl -d "test" localhost:8080
```

Via HTTP/2 (untested, my debian is old)
```shell
curl -d "test" --http2 localhost:8080
```

Via gRPC with [polyglot](https://github.com/grpc-ecosystem/polyglot)
```shell
echo '{"thing": "it"}' | java -jar polyglot.jar \
  --command=call \
  --endpoint localhost:8080 
  --full_method=doer.Doer/DoIt 
  --proto_discovery_root=doer 
  --use_tls=false
```
