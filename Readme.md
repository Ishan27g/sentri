#### start nats
```shell
docker run nats:alpine3.15
```

#### start listener
```shell
sentri
```

#### hook some commands
```shell
sentri go run database/main.go
sentri go run proxy/main.go
sentri go run auth/main.go
sentri go run cache/main.go
```

#### use as logrus hook
```go
package main

import "github.com/Ishan27g/sentri/hook"

func main(){
	
    w := hook.Hook("cmdName", os.Stdout)
    logrus.SetOutput(w)
	
}
```
