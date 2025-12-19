```md
$ brew install protobuf
grpc-connectRpc (main) $ go install github.com/bufbuild/buf/cmd/buf@v1.0.0-rc12     
grpc-connectRpc (main) $ go mod init rahulxf.com/rpc-learning                           
grpc-connectRpc (main) $ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
```

```md
buf config init
buf lint
buf generate
```

---

## DEMO 

```md
grpc-connectRpc (main) $ curl -X POST http://localhost:8080/auth.v1.AuthService/Register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@example.com",
    "password": "password123",
    "name": "Demo User"
  }'

{"userId":"6944eb986a92a9069778097b", "message":"User registered successfully"}%
```
```md
grpc-connectRpc (main) $ curl -X POST http://localhost:8080/auth.v1.AuthService/Login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@example.com",
    "password": "password123"
  }'

{"userId":"6944eb986a92a9069778097b", "accessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjk0NGViOTg2YTkyYTkwNjk3NzgwOTdiIiwid29ya3NwYWNlX2lkIjpbIjY5NDRlYjk4NmE5MmE5MDY5Nzc4MDk3YyJdLCJleHAiOjE3NjYxMjgwNDcsIm5iZiI6MTc2NjEyNDQ0NywiaWF0IjoxNzY2MTI0NDQ3fQ.N_cenxn_z3vXsY9_oQMhOM7vYUe5lTJ_hX3QxoTBwLk", "expiresAt":"1766128047"}%                                                                      
```

```md
grpc-connectRpc (main) $ export TOKEN="YOUR_TOKEN"
grpc-connectRpc (main) $ WORKSPACE_ID="6944eb986a92a9069778097c"
grpc-connectRpc (main) $ curl -X POST http://localhost:8080/todo.v1.TodoService/CreateTodo \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"workspace_id\": \"$WORKSPACE_ID\",
    \"title\": \"Learn ConnectRPC\",
    \"description\": \"Build a real backend\"
  }"

{"todo":{"id":"6944ec1c6a92a9069778097e", "workspaceId":"6944eb986a92a9069778097c", "title":"Learn ConnectRPC", "description":"Build a real backend", "createdBy":"6944eb986a92a9069778097b", "createdAt":"1766124572", "updatedAt":"1766124572"}}%                               grpc-connectRpc (main) $ TODO_ID="6944ec1c6a92a9069778097e"
```

```md
grpc-connectRpc (main) $ curl -X POST http://localhost:8080/todo.v1.TodoService/GetTodo \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"id\": \"$TODO_ID\",
    \"workspace_id\": \"$WORKSPACE_ID\"
  }"

{"todo":{"id":"6944ec1c6a92a9069778097e", "workspaceId":"6944eb986a92a9069778097c", "title":"Learn ConnectRPC", "description":"Build a real backend", "createdBy":"6944eb986a92a9069778097b", "createdAt":"1766124572", "updatedAt":"1766124572"}}%                               
```

```md
grpc-connectRpc (main) $ curl -X POST http://localhost:8080/todo.v1.TodoService/UpdateTodo \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"id\": \"$TODO_ID\",
    \"workspace_id\": \"$WORKSPACE_ID\",
    \"completed\": true
  }"

{"todo":{"id":"6944ec1c6a92a9069778097e", "workspaceId":"6944eb986a92a9069778097c", "title":"Learn ConnectRPC", "description":"Build a real backend", "completed":true, "createdBy":"6944eb986a92a9069778097b", "createdAt":"1766124572", "updatedAt":"1766124597"}}%             
```

```md
grpc-connectRpc (main) $ curl -X POST http://localhost:8080/todo.v1.TodoService/DeleteTodo \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"id\": \"$TODO_ID\",
    \"workspace_id\": \"$WORKSPACE_ID\"
  }"

{"success":true}%    
```

```md
grpc-connectRpc (main) $ curl -X POST http://localhost:8080/auth.v1.AuthService/Logout \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}'

{"success":true}%                                                                                                                        grpc-connectRpc (main) $
```