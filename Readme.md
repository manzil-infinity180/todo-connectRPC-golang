```md
$ brew install protobuf
grpc-connectRpc (main) $ go install github.com/bufbuild/buf/cmd/buf@v1.0.0-rc12     
grpc-connectRpc (main) $ go mod init rahulxf.com/rpc-learning                           
grpc-connectRpc (main) $ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
```

```md
buf config init
```

```md 
grpc-connectRpc (main) $ curl \
> --header "Content-Type: application/json" \
> --data '{"title":"Hello babu", "description":"teri to"}' \
> http://localhost:8080/todo.v1.TodoService/CreateTodo

{"todo":{"id":"generated-id","title":"Hello babu","description":"hello bhai kaise hoo","createdAt":"1765314586"}}
```