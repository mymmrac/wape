module github.com/mymmrac/wasm-gate

go 1.23.4

require (
	github.com/extism/go-pdk v1.1.0
	github.com/extism/go-sdk v1.6.1
	github.com/tetratelabs/wazero v1.8.2
)

replace github.com/extism/go-sdk => github.com/mymmrac/extism-go-sdk v0.0.0-20250109172127-8d6c4588106e

require (
	github.com/dylibso/observe-sdk/go v0.0.0-20240819160327-2d926c5d788a // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/ianlancetaylor/demangle v0.0.0-20240805132620-81f5be970eca // indirect
	github.com/tetratelabs/wabin v0.0.0-20230304001439-f6f874872834 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
