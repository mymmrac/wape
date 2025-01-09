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
	// ==== Envs ====
	// Defaults to none.

	// Envs sets an environment variables.
	Envs []string
	// EnvsMap sets an environment variables as a map.
	EnvsMap map[string]string
	// EnvsFile sets an environment variables from a file.
	EnvsFile string
	// EnvsFromHost pass thought environment variables from the host.
	EnvsFromHost bool

	// ==== Args ====
	// Defaults to none.

	// Args assigns command-line arguments.
	Args []string
	// ArgsFile assigns command-line arguments from a file.
	ArgsFile string
	// ArgsFromHost pass thought command-line arguments from the host.
	ArgsFromHost bool

	// ==== Stdin ====
	// Defaults to always return io.EOF.

	// Stdin configures where standard input (file descriptor 0) is read.
	Stdin io.Reader
	// StdinFile configures standard input (file descriptor 0) to read from a file.
	StdinFile string
	// StdinFromHost pass thought stdin from the host.
	StdinFromHost bool

	// ==== Stdout ====
	// Defaults to io.Discard.

	// Stdout configures where standard output (file descriptor 1) is written.
	Stdout io.Writer
	// StdoutFile configures standard output (file descriptor 1) to write to a file.
	StdoutFile string
	// StdoutFromHost pass thought stdout from the host.
	StdoutFromHost bool

	// ==== Stderr ====
	// Defaults to io.Discard.

	// Stderr configures where standard error (file descriptor 2) is written.
	Stderr io.Writer
	// StderrFile configures standard error (file descriptor 2) to write to a file.
	StderrFile string
	// StderrFromHost pass thought stderr from the host.
	StderrFromHost bool

	// ==== FS ====
	// Defaults to have no file access.

	// FS configures the filesystem.
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

	// ==== Random Source ====
	// Defaults to a deterministic source.

	// RandSource configures a source of random bytes.
	RandSource io.Reader
	// RandSourceFile configures a source of random bytes from a file.
	RandSourceFile string
	// RandSourceFromHost pass thought random source from the host.
	RandSourceFromHost bool

	// ==== Wall Time ====
	// Defaults to always return UTC 1 January 1970.

	// WallTime returns the current unix/epoch time, seconds since midnight UTC 1 January 1970,
	// with a nanosecond fraction.
	WallTime func() (sec int64, ns int32)
	// WallTimeFromHost pass thought wall time from the host.
	WallTimeFromHost bool

	// ==== Nano Time ====
	// Defaults to always return 0.

	// NanoTime returns nanoseconds since an arbitrary start point, used to measure elapsed time.
	// This is sometimes referred to as a tick or monotonic time.
	NanoTime func() int64
	// NanoTimeFromHost pass thought nano time from the host.
	NanoTimeFromHost bool

	// ==== Nano Sleep ====
	// Defaults to always return immediately.

	// NanoSleep puts the current goroutine to sleep for at least ns nanoseconds.
	NanoSleep func(ns int64)
	// NanoSleepFromHost pass thought nano sleep from the host.
	NanoSleepFromHost bool

	// ==== Start Functions ====
	// Defaults to "_start".

	// StartFunctions configures the functions to call after the module is instantiated.
	StartFunctions []string

	// ==== Memory ====

	// MemoryLimitPages overrides the maximum pages allowed per memory. Defaults to 65536, allowing 4GB total memory
	// per instance if the maximum is not encoded in a Wasm binary. Max is 65536 (2^16) pages or 4GB.
	MemoryLimitPages uint32

	// MemoryCapacityFromMax eagerly allocates max memory. Defaults to false, which means minimum memory is
	// allocated and any call to grow memory results in re-allocations.
	MemoryCapacityFromMax bool

	// ==== Execution Time ====

	// MaxExecutionDuration limits the maximum execution time. Defaults to 0 (no limit).
	MaxExecutionDuration time.Duration

	// ==== Debug ====

	// DebugInfoEnabled toggles DWARF-based stack traces in the face of runtime errors. Defaults to false.
	DebugInfoEnabled bool

	// ==== Network ====

	// TODO: Add network access

	// ==== WASI ====

	// DisableWASI disables WASI Preview 1 support. Defaults to false.
	DisableWASIP1 bool

	// ==== Wazero ===

	// ModuleConfig configures the WASM module.
	ModuleConfig wazero.ModuleConfig
	// RuntimeConfig configures the WASM runtime.
	RuntimeConfig wazero.RuntimeConfig

	// ==== Extism ====

	// Manifest is the plugin manifest that can be provided instead of configuration above.
	Manifest *extism.Manifest
	// PluginConfig is the plugin configuration that can be provided instead of configuration above.
	PluginConfig *extism.PluginConfig

	// ==== Host Functions ====

	// HostFunctions host functions available to the guest WASM module.
	HostFunctions []extism.HostFunction
}

// NewEnvironment returns a new environment.
func NewEnvironment() *Environment {
	return &Environment{}
}

// MakeModuleConfig returns the module configuration based on the environment.
func (e *Environment) MakeModuleConfig() wazero.ModuleConfig {
	if e.ModuleConfig != nil {
		return e.ModuleConfig
	}
	return wazero.NewModuleConfig()
}

// MakeRuntimeConfig returns the runtime configuration based on the environment.
func (e *Environment) MakeRuntimeConfig() wazero.RuntimeConfig {
	if e.RuntimeConfig != nil {
		return e.RuntimeConfig
	}
	return wazero.NewRuntimeConfig()
}

// MakeManifest returns the manifest based on the environment.
func (e *Environment) MakeManifest() extism.Manifest {
	if e.Manifest != nil {
		return *e.Manifest
	}
	return extism.Manifest{}
}

// MakePluginConfig returns the plugin configuration based on the environment.
func (e *Environment) MakePluginConfig() extism.PluginConfig {
	if e.PluginConfig != nil {
		return *e.PluginConfig
	}
	return extism.PluginConfig{}
}
