package wasmgate

import (
	"io"
	"io/fs"
	"time"

	extism "github.com/extism/go-sdk"
	"github.com/tetratelabs/wazero"
)

// Environment configures the behavior of WASM module.
//
// Note: Many environment configurations override each other, for example providing ModuleConfig will override all
// other configurations related to envs, args, FS, etc. Be aware that you may end up with unexpected behavior.
//
// TODO: Add json/yaml/toml tags
type Environment struct {
	// Envs sets an environment variables. Defaults to none.
	Envs []string
	// EnvsFile sets an environment variables from a file.
	EnvsFile string
	// EnvsFromHost pass thought environment variables from the host.
	EnvsFromHost bool

	// Args assigns command-line arguments. Defaults to none.
	Args []string
	// ArgsFile assigns command-line arguments from a file.
	ArgsFile string
	// ArgsFromHost pass thought command-line arguments from the host.
	ArgsFromHost bool

	// Stdin configures where standard input (file descriptor 0) is read. Defaults to return io.EOF.
	Stdin io.Reader
	// StdinFile configures standard input (file descriptor 0) to read from a file.
	StdinFile string
	// StdinFromHost pass thought stdin from the host.
	StdinFromHost bool

	// Stdout configures where standard output (file descriptor 1) is written. Defaults to io.Discard.
	Stdout io.Writer
	// StdoutFile configures standard output (file descriptor 1) to write to a file.
	StdoutFile string
	// StdoutFromHost pass thought stdout from the host.
	StdoutFromHost bool

	// Stderr configures where standard error (file descriptor 2) is written. Defaults to io.Discard.
	Stderr io.Writer
	// StderrFile configures standard error (file descriptor 2) to write to a file.
	StderrFile string
	// StderrFromHost pass thought stderr from the host.
	StderrFromHost bool

	// FS configures the filesystem. Defaults to have no file access is allowed.
	FS fs.FS
	// FSFile configures the filesystem mount points.
	FSConfig wazero.FSConfig
	// FSDir configures the filesystem as a root directory.
	FSDir string
	// FSAllowedPaths configures the allowed filesystem paths that will be mapped in WASM module.
	FSAllowedPaths map[string]string
	// FSFromHost pass thought filesystem from the host.
	FSFromHost bool

	// TODO: Add file permissions + create/remove files

	// RandSource configures a source of random bytes. Defaults to a deterministic source.
	RandSource io.Reader
	// RandSourceFile configures a source of random bytes from a file.
	RandSourceFile string
	// RandSourceFromHost pass thought random source from the host.
	RandSourceFromHost bool

	// WallTime returns the current unix/epoch time, seconds since midnight UTC 1 January 1970,
	// with a nanosecond fraction. Defaults to always return UTC 1 January 1970.
	WallTime func() (sec int64, ns int32)
	// WallTimeFromHost pass thought wall time from the host.
	WallTimeFromHost bool

	// NanoTime returns nanoseconds since an arbitrary start point, used to measure elapsed time.
	// This is sometimes referred to as a tick or monotonic time. Defaults to always return 0.
	NanoTime func() int64
	// NanoTimeFromHost pass thought nano time from the host.
	NanoTimeFromHost bool

	// NanoSleep puts the current goroutine to sleep for at least ns nanoseconds. Defaults to return immediately.
	NanoSleep func(ns int64)
	// NanoSleepFromHost pass thought nano sleep from the host.
	NanoSleepFromHost bool

	// StartFunctions configures the functions to call after the module is instantiated. Defaults to "_start".
	StartFunctions []string

	// MemoryLimitPages overrides the maximum pages allowed per memory. Defaults to 65536, allowing 4GB total memory
	// per instance if the maximum is not encoded in a Wasm binary. Max is 65536 (2^16) pages or 4GB.
	MemoryLimitPages uint32

	// MemoryCapacityFromMax eagerly allocates max memory. Defaults to false, which means minimum memory is
	// allocated and any call to grow memory results in re-allocations.
	MemoryCapacityFromMax bool

	// MaxExecutionDuration limits the maximum execution time. Defaults to 0 (no limit).
	MaxExecutionDuration time.Duration

	// DebugInfoEnabled toggles DWARF-based stack traces in the face of runtime errors. Defaults to false.
	DebugInfoEnabled bool

	// TODO: Add network access

	// DisableWASI disables WASI Preview 1 support. Defaults to false.
	DisableWASIP1 bool

	// ModuleConfig configures the WASM module.
	ModuleConfig wazero.ModuleConfig
	// RuntimeConfig configures the WASM runtime.
	RuntimeConfig wazero.RuntimeConfig

	// Manifest is the plugin manifest that can be provided instead of configuration above.
	Manifest extism.Manifest
	// PluginConfig is the plugin configuration that can be provided instead of configuration above.
	PluginConfig extism.PluginConfig

	// HostFunctions host functions available to the guest WASM module.
	HostFunctions []extism.HostFunction
}
