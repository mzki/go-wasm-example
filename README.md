# go-wasm-example

This repository contains example of how to run wasm binary compiled by Go and TinyGo.
These exmaples are tested only on a web browser, not WASI compatiable envrionment.

The example shows:
* Exporting Go function to JS side (which is same as TinyGo exmaple). See `multiple` function in the code.
* Call async JS function from Go side. See `addAsync` function in the code.
* Call async JS function which takes a callback function as argument from Go side. See `addPromiss` function in the code.

Since TinyGo `0.26.0`, `runtime.timer` is implemented for wasm platform, Go code using timer feature, including `context.WithTimeout`, can be compiled. Older TinyGo version `<0.26.0` will not work this code.

## Usage

* Build wasm binary

```bash
# using tiny go
bash scripts/build-tiny-wasm-on-docker
# using go
bash scripts/build-wasm-on-docker
```

* Serve compiled binary and html

```bash
go run server.go --dir tiny-html
# if you use go version instead of tinygo version, use below
go run server.go --dir html
```

* Access `localhost:8080` from web broser.

* Result are shown on Debug console (The screen still blank page since no content on html.)
```
  GEThttp://localhost:8080/favicon.ico
[HTTP/1.1 404 Not Found 0ms]

multiplied two numbers: 15 
adding promiss two numbers: 5
adding async two numbers: 5
```
