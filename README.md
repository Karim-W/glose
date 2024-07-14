# Glose

Global Golang Close handler to handle graceful shutdown of your application in
a global way without having to pass around a context.

## Usage

### Shutdown on Panics

```go
package main

import (
    "net/http"
    "time"

    "github.com/karim-w/glose"
)

type watchdog struct {
    // Your watchdog struct
}

func (w *watchdog) Close() {
    // Clean up your resources here
}

type MyServer struct {
    // Your server struct
}

func (s *MyServer) Close() {
    // Clean up your resources here
}

func (s *MyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Your handler code here
}

func function_that_panics() {
    glose.Panik("I'm panicking!")
}

func main() {
    server := &MyServer{}

    // manually register entities to the global glose handler
    glose.Register(&watchdog{}, server)

    go func() {
        time.Sleep(time.Second * 2)
        function_that_panics()
    }()


    http.Handle("/", server)
}
```

### Shutdown on Signals

```go
package main

import (
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/karim-w/glose"
)


type watchdog struct {
    // Your watchdog struct
}

func (w *watchdog) Close() {
    // Clean up your resources here
}

type MyServer struct {
    // Your server struct
}

func (s *MyServer) Close() {
    // Clean up your resources here
}

func (s *MyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Your handler code here
}

func function_that_panics() {
    glose.Panik("I'm panicking!")
}

func main() {
    server := &MyServer{}

    go func() {
        http.Handle("/", server)
        http.ListenAndServe(":8080", nil)
    }()

    go func() {
        time.Sleep(time.Second * 2)
        function_that_panics()
    }()


    // closes on graceful shutdown
    glose.Watch(
        &watchdog{},
        server,
    )
}
```

## License

BSD 3-Clause License

## Author

karim-w

## Contributing

Pull requests are welcome. For major changes, please open an issue first to
discuss what you would like to change.
