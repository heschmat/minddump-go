
## fundamentals
```sh
go version

go mod init github.com/heschmat/minddump-go

```
When there's a valid `go.mod` file in the root of your project directory, your project is a **module**.


### web application basics

Handlers are responsible for executing the application logic & writing HTTP responses. 

The router/servemux stores a mapping between the URL & handlers. It dispaches the reques to the matching handler.

In Go we can simply establish a web server & listen for incoming requests as part of the application itself.

`* http.Request` parameter is a pointer to a struct which holds information about the current request.

N.B. Go's servemux treats the route patter "/" like a catch-all.

`go run` is a shortcut to combile the code, create an executable binary in `/tmp` & then run it.
```sh
go run .
go run main.go
go run `module-path`
```

### 

subtree path patter (when a route pattern ends with a trailing slash: "/" or "/static")

`r.PathValue()` to get the value for a wildcard segment: `r.PathValue("itemID")`. It returns a string.

**The most specific pattern wins**. But avoid it.

### customizing responses

If you want t send a non-200 status code, you must call `w.WriteHeader()` before any call to `w.Write()`

#### writing response bodies
In practice, it's common to pass the `http.ResponseWriter` value to another function that writes the resposne `w` for us.

**...because the `http.ResponseWriter` value in the handlers has a `.Write()` method, it satisfies the `io.Write` interface.**
```go
// instead of 
w.Write([]byte("Bonjour"))

// either of:
io.WriteString(w, "Bonjour")
fmt.Fprint(w, "Bonjour")

```

### serving static files

```md
URL:  /static/css/main.css
        │
        ▼
ServeMux route match (GET /static/)
        │
        ▼
StripPrefix("/static")
        │
        ▼
New path: /css/main.css
        │
        ▼
FileServer root: ./ui/static/
        │
        ▼
Disk path: ./ui/static/css/main.css
        │
        ▼
Serve file (or 404)

```


So why not serve `http.FileServer(http.Dir("./ui/"))`?
Because **FileServer** exposes everything reachable under that root via HTTP.
Not just what you link to.
Anything a user guesses becomes accessible.

## configuration & error handling

❌ `addr := os.GetEnv("ADDR_")`

Use flags:
```go
// `flag.String()` returns a pointer to the flag value (not the value itself)
addr := flag.String("addr", ":4000", "HTTP network address")
flag.Parse()

// pass the dereferenced `addr` pointer
log.Printf("starting the server on %s", *addr)

// go run ./cmd/web -addr=$ADDR_
```
To list all the available command-line flags: `go run ./cmd/web -help` 

### Dependency injection
Q: how to make any dependency - structurred logger, db connection pool, centralized error handler ...- available to the handlers?

A: inject dependencies into the handlers

As all our handlers are in the same package, we can inject dependencies like so:
1- put the dpendencies into a custom **application struct**
2- define the handlers as methods agains the struct (which holds the application-wide dependencies)


## MiSC

byte slice 
---
`cat /etc/services | grep http-alt`
---

### funcs
```go
http.NotFound(w, r)
```

### log

```go
log.Println("starting the server on :4000")
log.Printf("starting the server on %s", *addr)

# -----------
err := http.ListenAndServe(":4000", nil)
log.Fatal(err)
```

#### structured logging 

```go
// args:
// 1: write destination for the log entries
// 2: a pointer to a `slog.HandlerOptions` to customize the behavior (nil for defaul)
loggerHandler := slog.NewTextHandler(os.Stdout, nil)
// create the structured logger
logger := slog.New(loggerHandler)
```
During development, the logs are displayed in the terminal (the standard output.)

In staging/prod we can redirect the standard out stream to an on-disk file: `go run ./cmd/web >> /tmp/minddump.log`

N.B. custom loggers created by `slog.New()` are concurrently-safe. 

### fmt
```go
fmt.Fprint(w, "Bonjour")

fmt.Sprintf("item %d...", id)



// instead of:
msg := fmt.Sprintf("snippet %d...", id)
w.Write([]byte(msg))
// simply
fmt.Fprintf(w, "snippet %d...", id)
```

#### err
```go
if err != nil {
  log.Println(err.Error())
  http.Error(w, "Internal Server Error", http.StatusInternalServerError)
  return
}
```

### curl

```sh
curl -i localhost:4000/

curl --head localhost:4000/

curl -i -d "" localhost:4000/


curl -i localhost:4000/snippet/create
curl -i -d "" localhost:4000/snippet/create
curl -i -X DELETE localhost:4000/snippet/create



```

### statusCode
```md
405 Method Not Allowed

```

## structure

`cmd/` contains the **application-specific code**.
```sh
mkdir -p cmd/web internal ui/html ui/static

# to automate some administrative tasks in the future
# mkdir -p cmd/cli

```

### cmd/web
```md
main.go
handlers.go
```

### internal
contains the ancillary non-application-specific code.

### ui
ui/html/pages
    home.tmpl
    

ui/static