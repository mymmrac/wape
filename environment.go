package wape

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	extism "github.com/extism/go-sdk"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/sys"

	wio "github.com/mymmrac/wape/host/io"
	wnet "github.com/mymmrac/wape/host/net"
)

// Environment configures the behavior of WASM module.
//
// Note: Many environment configurations override each other, for example providing ModuleConfig will override all
// other configurations related to envs, args, FS, etc. Be aware that you may end up with unexpected behavior.
type Environment struct {
	// ==== Envs ====
	// Defaults to none.

	// Envs sets an environment variables.
	Envs []string `json:"envs,omitempty" yaml:"envs,omitempty" toml:"envs,omitempty"`
	// EnvsMap sets an environment variables as a map.
	EnvsMap map[string]string `json:"envsMap,omitempty" yaml:"envsMap,omitempty" toml:"envsMap,omitempty"`
	// EnvsFile sets an environment variables from a file.
	EnvsFile string `json:"envsFile,omitempty" yaml:"envsFile,omitempty" toml:"envsFile,omitempty"`
	// EnvsFromHost pass thought environment variables from the host.
	EnvsFromHost bool `json:"envsFromHost,omitempty" yaml:"envsFromHost,omitempty" toml:"envsFromHost,omitempty"`

	// ==== Args ====
	// Defaults to none.

	// Args assigns command-line arguments.
	Args []string `json:"args,omitempty" yaml:"args,omitempty" toml:"args,omitempty"`
	// ArgsFile assigns command-line arguments from a file.
	ArgsFile string `json:"argsFile,omitempty" yaml:"argsFile,omitempty" toml:"argsFile,omitempty"`
	// ArgsFromHost pass thought command-line arguments from the host.
	ArgsFromHost bool `json:"argsFromHost,omitempty" yaml:"argsFromHost,omitempty" toml:"argsFromHost,omitempty"`

	// ==== Stdin ====
	// Defaults to always return io.EOF.

	// Stdin configures where standard input (file descriptor 0) is read.
	Stdin io.Reader `json:"-" yaml:"-" toml:"-"`
	// StdinFile configures standard input (file descriptor 0) to read from a file.
	StdinFile string `json:"stdinFile,omitempty" yaml:"stdinFile,omitempty" toml:"stdinFile,omitempty"`
	// StdinFromHost pass thought stdin from the host.
	StdinFromHost bool `json:"stdinFromHost,omitempty" yaml:"stdinFromHost,omitempty" toml:"stdinFromHost,omitempty"`

	// ==== Stdout ====
	// Defaults to io.Discard.

	// Stdout configures where standard output (file descriptor 1) is written.
	Stdout io.Writer `json:"-" yaml:"-" toml:"-"`
	// StdoutFile configures standard output (file descriptor 1) to write to a file.
	StdoutFile string `json:"stdoutFile,omitempty" yaml:"stdoutFile,omitempty" toml:"stdoutFile,omitempty"`
	// StdoutFromHost pass thought stdout from the host.
	StdoutFromHost bool `json:"stdoutFromHost,omitempty" yaml:"stdoutFromHost,omitempty" toml:"stdoutFromHost,omitempty"`

	// ==== Stderr ====
	// Defaults to io.Discard.

	// Stderr configures where standard error (file descriptor 2) is written.
	Stderr io.Writer `json:"-" yaml:"-" toml:"-"`
	// StderrFile configures standard error (file descriptor 2) to write to a file.
	StderrFile string `json:"stderrFile,omitempty" yaml:"stderrFile,omitempty" toml:"stderrFile,omitempty"`
	// StderrFromHost pass thought stderr from the host.
	StderrFromHost bool `json:"stderrFromHost,omitempty" yaml:"stderrFromHost,omitempty" toml:"stderrFromHost,omitempty"`

	// ==== FS ====
	// Defaults to have no file access.

	// FS configures the filesystem.
	FS fs.FS `json:"-" yaml:"-" toml:"-"`
	// FSFile configures the filesystem mount points.
	FSConfig wazero.FSConfig `json:"-" yaml:"-" toml:"-"`
	// FSDir configures the filesystem as a root directory.
	FSDir string `json:"fsDir,omitempty" yaml:"fsDir,omitempty" toml:"fsDir,omitempty"`
	// FSAllowedPaths configures the allowed filesystem paths that will be mapped in WASM module.
	// Host paths prefixed with "ro:" are marked as read-only.
	FSAllowedPaths map[string]string `json:"fsAllowedPaths,omitempty" yaml:"fsAllowedPaths,omitempty" toml:"fsAllowedPaths,omitempty"`
	// FSFromHost pass thought filesystem from the host.
	FSFromHost bool `json:"fsFromHost,omitempty" yaml:"fsFromHost,omitempty" toml:"fsFromHost,omitempty"`

	// ==== Random Source ====
	// Defaults to a deterministic source.

	// RandSource configures a source of random bytes.
	RandSource io.Reader `json:"-" yaml:"-" toml:"-"`
	// RandSourceFile configures a source of random bytes from a file.
	RandSourceFile string `json:"randSourceFile,omitempty" yaml:"randSourceFile,omitempty" toml:"randSourceFile,omitempty"`
	// RandSourceFromHost pass thought random source from the host.
	RandSourceFromHost bool `json:"randSourceFromHost,omitempty" yaml:"randSourceFromHost,omitempty" toml:"randSourceFromHost,omitempty"`

	// ==== Wall Time ====
	// Defaults to always return UTC 1 January 1970.

	// WallTime returns the current unix/epoch time, seconds since midnight UTC 1 January 1970,
	// with a nanosecond fraction.
	WallTime sys.Walltime `json:"-" yaml:"-" toml:"-"`
	// WallTimeClockResolution configures the resolution of the wall clock, defaults to 1ns.
	WallTimeClockResolution sys.ClockResolution `json:"-" yaml:"-" toml:"-"`
	// WallTimeFromHost pass thought wall time from the host.
	WallTimeFromHost bool `json:"wallTimeFromHost,omitempty" yaml:"wallTimeFromHost,omitempty" toml:"wallTimeFromHost,omitempty"`

	// ==== Nano Time ====
	// Defaults to always return 0.

	// NanoTime returns nanoseconds since an arbitrary start point, used to measure elapsed time.
	// This is sometimes referred to as a tick or monotonic time.
	NanoTime sys.Nanotime `json:"-" yaml:"-" toml:"-"`
	// NanoTimeClockResolution configures the resolution of the nano clock, defaults to 1ns.
	NanoTimeClockResolution sys.ClockResolution `json:"-" yaml:"-" toml:"-"`
	// NanoTimeFromHost pass thought nano time from the host.
	NanoTimeFromHost bool `json:"nanoTimeFromHost,omitempty" yaml:"nanoTimeFromHost,omitempty" toml:"nanoTimeFromHost,omitempty"`

	// ==== Nano Sleep ====
	// Defaults to always return immediately.

	// NanoSleep puts the current goroutine to sleep for at least ns nanoseconds.
	NanoSleep sys.Nanosleep `json:"-" yaml:"-" toml:"-"`
	// NanoSleepFromHost pass thought nano sleep from the host.
	NanoSleepFromHost bool `json:"nanoSleepFromHost,omitempty" yaml:"nanoSleepFromHost,omitempty" toml:"nanoSleepFromHost,omitempty"`

	// ==== Start Functions ====
	// Defaults to none.

	// StartFunctions configures the functions to call after the module is instantiated.
	StartFunctions []string `json:"startFunctions,omitempty" yaml:"startFunctions,omitempty" toml:"startFunctions,omitempty"`

	// ==== Memory ====

	// MemoryLimitPages overrides the maximum pages allowed per memory. Defaults to 65536, allowing 4GB total memory
	// per instance if the maximum is not encoded in a WASM binary. Max is 65536 (2^16) pages or 4GB.
	MemoryLimitPages uint32 `json:"memoryLimitPages,omitempty" yaml:"memoryLimitPages,omitempty" toml:"memoryLimitPages,omitempty"`

	// MemoryCapacityFromMax eagerly allocates max memory. Defaults to false, which means minimum memory is
	// allocated and any call to grow memory results in re-allocations.
	MemoryCapacityFromMax bool `json:"memoryCapacityFromMax,omitempty" yaml:"memoryCapacityFromMax,omitempty" toml:"memoryCapacityFromMax,omitempty"`

	// ==== Execution Time ====

	// MaxExecutionDuration limits the maximum function execution time.
	// Rounded to milliseconds and has a minimum of 1ms. Defaults to 0 (no limit).
	MaxExecutionDuration time.Duration `json:"maxExecutionDuration,omitempty" yaml:"maxExecutionDuration,omitempty" toml:"maxExecutionDuration,omitempty"`

	// ==== Debug ====

	// DebugInfoEnabled toggles DWARF-based stack traces in the face of runtime errors. Defaults to false.
	DebugInfoEnabled bool `json:"debugInfoEnabled,omitempty" yaml:"debugInfoEnabled,omitempty" toml:"debugInfoEnabled,omitempty"`

	// ExtismDebugEnvAllowed allows use of EXTISM_ENABLE_WASI_OUTPUT environment variable. Defaults to false (unsets
	// this variable, so it will not work for all WASM modules).
	ExtismDebugEnvAllowed bool `json:"extismDebugEnvAllowed,omitempty" yaml:"extismDebugEnvAllowed,omitempty" toml:"extismDebugEnvAllowed,omitempty"`

	// ==== Network ====

	// NetworkEnabled toggles network access. Defaults to false.
	NetworkEnabled bool `json:"networkEnabled,omitempty" yaml:"networkEnabled,omitempty" toml:"networkEnabled,omitempty"`

	// NetworksAllowed configures the allowed network protocols. See [net.Dial] for allowed protocols. Defaults to none.
	NetworksAllowed []string `json:"networksAllowed,omitempty" yaml:"networksAllowed,omitempty" toml:"networksAllowed,omitempty"`
	// NetworksAllowAll allows all network protocols. Defaults to false.
	NetworksAllowAll bool `json:"networksAllowAll,omitempty" yaml:"networksAllowAll,omitempty" toml:"networksAllowAll,omitempty"`

	// NetworkAddressesAllowed configures the allowed network addresses. Defaults to none.
	NetworkAddressesAllowed []string `json:"networkAddressesAllowed,omitempty" yaml:"networkAddressesAllowed,omitempty" toml:"networkAddressesAllowed,omitempty"`
	// NetworkAddressesAllowAll allows all network addresses. Defaults to false.
	NetworkAddressesAllowAll bool `json:"networkAddressesAllowAll,omitempty" yaml:"networkAddressesAllowAll,omitempty" toml:"networkAddressesAllowAll,omitempty"`

	// ==== WASI ====

	// DisableWASI disables WASI Preview 1 support. Defaults to false.
	DisableWASI bool `json:"disableWASI,omitempty" yaml:"disableWASI,omitempty" toml:"disableWASI,omitempty"`

	// ==== Wazero ===

	// CompilationCache configures the compilation cache.
	CompilationCache wazero.CompilationCache `json:"-" yaml:"-" toml:"-"`
	// CompilationCacheDir configures the compilation cache directory.
	CompilationCacheDir string `json:"compilationCacheDir,omitempty" yaml:"compilationCacheDir,omitempty" toml:"compilationCacheDir,omitempty"`

	// CustomSectionsEnabled configures whether to enable custom sections.
	CustomSectionsEnabled bool `json:"customSectionsEnabled,omitempty" yaml:"customSectionsEnabled,omitempty" toml:"customSectionsEnabled,omitempty"`

	// ModuleConfig configures the WASM module.
	ModuleConfig wazero.ModuleConfig `json:"-" yaml:"-" toml:"-"`
	// RuntimeConfig configures the WASM runtime.
	RuntimeConfig wazero.RuntimeConfig `json:"-" yaml:"-" toml:"-"`

	// ==== Extism ====

	// Manifest is the plugin manifest that can be provided instead of configuration above.
	Manifest *extism.Manifest `json:"manifest,omitempty" yaml:"manifest,omitempty" toml:"manifest,omitempty"`
	// PluginConfig is the plugin configuration that can be provided instead of configuration above.
	PluginConfig *extism.PluginConfig `json:"-" yaml:"-" toml:"-"`
	// PluginInstanceConfig is the plugin instance configuration that can be provided instead of configuration above.
	PluginInstanceConfig *extism.PluginInstanceConfig `json:"-" yaml:"-" toml:"-"`

	// ==== Host Functions ====

	// HostFunctions host functions available to the guest WASM module.
	HostFunctions []extism.HostFunction `json:"-" yaml:"-" toml:"-"`

	// ==== WASM Modules ====

	// Modules configures the WASM modules.
	Modules []ModuleData `json:"modules,omitempty" yaml:"modules,omitempty" toml:"modules,omitempty"`
}

// ModuleData represents WASM module with name and data.
type ModuleData struct {
	// Name of the WASM module.
	Name string `json:"name" yaml:"name" toml:"name"`

	// Source of WASM module.
	Data []byte `json:"data,omitempty" yaml:"data,omitempty" toml:"data,omitempty"`

	// File path to read WASM module.
	File string `json:"file,omitempty" yaml:"file,omitempty" toml:"file,omitempty"`

	// Url to download WASM module.
	Url string `json:"url,omitempty" yaml:"url,omitempty" toml:"url,omitempty"`
	// Method to download WASM module.
	// Defaults to "GET".
	HttpMethod string `json:"httpMethod,omitempty" yaml:"httpMethod,omitempty" toml:"httpMethod,omitempty"`
	// Headers to download WASM module.
	HttpHeaders map[string]string `json:"httpHeaders,omitempty" yaml:"httpHeaders,omitempty" toml:"httpHeaders,omitempty"`

	// SHA256 hash of WASM module to validate it.
	Hash string `json:"hash,omitempty" yaml:"hash,omitempty" toml:"hash,omitempty"`
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

	cfg := wazero.NewModuleConfig()

	switch {
	case e.EnvsFromHost:
		for _, env := range os.Environ() {
			key, value, _ := strings.Cut(env, "=")
			cfg = cfg.WithEnv(key, value)
		}
	case e.EnvsFile != "":
		envs, err := os.ReadFile(e.EnvsFile)
		if err == nil {
			for env := range strings.Lines(string(envs)) {
				env = strings.TrimSpace(env)
				if env == "" {
					continue
				}
				key, value, _ := strings.Cut(env, "=")
				cfg = cfg.WithEnv(key, value)
			}
		}
	case len(e.EnvsMap) > 0:
		for key, value := range e.EnvsMap {
			cfg = cfg.WithEnv(key, value)
		}
	case len(e.Envs) > 0:
		for _, env := range e.Envs {
			key, value, _ := strings.Cut(env, "=")
			cfg = cfg.WithEnv(key, value)
		}
	}

	switch {
	case e.ArgsFromHost:
		cfg = cfg.WithArgs(os.Args...)
	case e.ArgsFile != "":
		args, err := os.ReadFile(e.ArgsFile)
		if err == nil {
			cfg = cfg.WithArgs(strings.Fields(string(args))...)
		}
	case len(e.Args) > 0:
		cfg = cfg.WithArgs(e.Args...)
	}

	switch {
	case e.StdinFromHost:
		cfg = cfg.WithStdin(os.Stdin)
	case e.StdinFile != "":
		stdin, err := os.Open(e.StdinFile)
		if err == nil {
			cfg = cfg.WithStdin(stdin)
		}
	case e.Stdin != nil:
		cfg = cfg.WithStdin(e.Stdin)
	}

	switch {
	case e.StdoutFromHost:
		cfg = cfg.WithStdout(os.Stdout)
	case e.StdoutFile != "":
		stdout, err := os.OpenFile(e.StdoutFile, os.O_WRONLY|os.O_CREATE, 0644)
		if err == nil {
			cfg = cfg.WithStdout(stdout)
		}
	case e.Stdout != nil:
		cfg = cfg.WithStdout(e.Stdout)
	}

	switch {
	case e.StderrFromHost:
		cfg = cfg.WithStderr(os.Stderr)
	case e.StderrFile != "":
		stderr, err := os.OpenFile(e.StderrFile, os.O_WRONLY|os.O_CREATE, 0644)
		if err == nil {
			cfg = cfg.WithStderr(stderr)
		}
	case e.Stderr != nil:
		cfg = cfg.WithStderr(e.Stderr)
	}

	switch {
	case e.FSFromHost:
		if runtime.GOOS == "windows" {
			wd, err := os.Getwd()
			if err != nil {
				cfg = cfg.WithFS(os.DirFS("C:\\"))
			} else {
				cfg = cfg.WithFS(os.DirFS(filepath.VolumeName(wd) + "\\"))
			}
		} else {
			cfg = cfg.WithFS(os.DirFS("/"))
		}
	case len(e.FSAllowedPaths) > 0:
		fsCfg := wazero.NewFSConfig()
		for host, guest := range e.FSAllowedPaths {
			if strings.HasPrefix(host, "ro:") {
				fsCfg = fsCfg.WithReadOnlyDirMount(strings.TrimPrefix(host, "ro:"), guest)
			} else {
				fsCfg = fsCfg.WithDirMount(host, guest)
			}
		}
		cfg = cfg.WithFSConfig(fsCfg)
	case e.FSDir != "":
		root, err := os.OpenRoot(e.FSDir)
		if err == nil {
			cfg = cfg.WithFS(root.FS())
		}
	case e.FSConfig != nil:
		cfg = cfg.WithFSConfig(e.FSConfig)
	case e.FS != nil:
		cfg = cfg.WithFS(e.FS)
	}

	switch {
	case e.RandSourceFromHost:
		cfg = cfg.WithRandSource(rand.Reader)
	case e.RandSourceFile != "":
		randSource, err := os.Open(e.RandSourceFile)
		if err == nil {
			cfg = cfg.WithRandSource(randSource)
		}
	case e.RandSource != nil:
		cfg = cfg.WithRandSource(e.RandSource)
	}

	switch {
	case e.WallTimeFromHost:
		cfg = cfg.WithSysWalltime()
	case e.WallTime != nil:
		if e.WallTimeClockResolution == 0 {
			e.WallTimeClockResolution = 1 // 1ns
		}
		cfg = cfg.WithWalltime(e.WallTime, e.WallTimeClockResolution)
	}

	switch {
	case e.NanoTimeFromHost:
		cfg = cfg.WithSysNanotime()
	case e.NanoTime != nil:
		if e.NanoTimeClockResolution == 0 {
			e.NanoTimeClockResolution = 1 // 1ns
		}
		cfg = cfg.WithNanotime(e.NanoTime, e.NanoTimeClockResolution)
	}

	switch {
	case e.NanoSleepFromHost:
		cfg = cfg.WithSysNanosleep()
	case e.NanoSleep != nil:
		cfg = cfg.WithNanosleep(e.NanoSleep)
	}

	cfg.WithStartFunctions(e.StartFunctions...)

	return cfg
}

// MakePluginInstanceConfig returns the plugin instance configuration based on the environment.
func (e *Environment) MakePluginInstanceConfig() extism.PluginInstanceConfig {
	if e.PluginInstanceConfig != nil {
		return *e.PluginInstanceConfig
	}

	return extism.PluginInstanceConfig{
		ModuleConfig: e.MakeModuleConfig(),
	}
}

// MakeRuntimeConfig returns the runtime configuration based on the environment.
func (e *Environment) MakeRuntimeConfig() wazero.RuntimeConfig {
	if e.RuntimeConfig != nil {
		return e.RuntimeConfig
	}

	cfg := wazero.NewRuntimeConfig()

	cfg = cfg.WithDebugInfoEnabled(e.DebugInfoEnabled)

	if e.MemoryLimitPages != 0 {
		cfg = cfg.WithMemoryLimitPages(e.MemoryLimitPages)
	}

	cfg = cfg.WithMemoryCapacityFromMax(e.MemoryCapacityFromMax)

	switch {
	case e.CompilationCacheDir != "":
		cache, err := wazero.NewCompilationCacheWithDir(e.CompilationCacheDir)
		if err == nil {
			cfg = cfg.WithCompilationCache(cache)
		}
	case e.CompilationCache != nil:
		cfg = cfg.WithCompilationCache(e.CompilationCache)
	}

	cfg = cfg.WithCustomSections(e.CustomSectionsEnabled)

	if e.MaxExecutionDuration > 0 {
		cfg = cfg.WithCloseOnContextDone(true)
	}

	if !e.ExtismDebugEnvAllowed {
		const env = "EXTISM_ENABLE_WASI_OUTPUT"
		err := os.Unsetenv(env)
		if err != nil {
			panic(fmt.Errorf("unset %q environment variable: %w", env, err))
		}
	}

	return cfg
}

// MakeManifest returns the manifest based on the environment.
func (e *Environment) MakeManifest() extism.Manifest {
	if e.Manifest != nil {
		return *e.Manifest
	}

	manifest := extism.Manifest{}

	if e.MaxExecutionDuration > 0 {
		manifest.Timeout = max(uint64(e.MaxExecutionDuration.Round(time.Millisecond).Milliseconds()), 1)
	}

	for _, module := range e.Modules {
		var wasm extism.Wasm

		switch {
		case module.Data != nil:
			wasm = &extism.WasmData{
				Data: module.Data,
				Hash: module.Hash,
				Name: module.Name,
			}
		case module.File != "":
			wasm = &extism.WasmFile{
				Path: module.File,
				Hash: module.Hash,
				Name: module.Name,
			}
		case module.Url != "":
			wasm = &extism.WasmUrl{
				Url:     module.Url,
				Hash:    module.Hash,
				Headers: module.HttpHeaders,
				Name:    module.Name,
				Method:  module.HttpMethod,
			}
		default:
			continue
		}

		manifest.Wasm = append(manifest.Wasm, wasm)
	}

	return manifest
}

// MakePluginConfig returns the plugin configuration based on the environment.
func (e *Environment) MakePluginConfig() extism.PluginConfig {
	if e.PluginConfig != nil {
		return *e.PluginConfig
	}
	return extism.PluginConfig{
		RuntimeConfig: e.MakeRuntimeConfig(),
		EnableWasi:    !e.DisableWASI,
		ModuleConfig:  e.MakeModuleConfig(),
	}
}

// MakeHostFunctions returns the host functions based on the environment.
func (e *Environment) MakeHostFunctions() []extism.HostFunction {
	functions := make([]extism.HostFunction, 0, len(e.HostFunctions))

	if e.NetworkEnabled {
		functions = append(functions, wio.Ready())
	}

	if e.NetworkEnabled {
		functions = append(functions, wnet.LookupHost())
		functions = append(functions, wnet.Dial(wnet.DialConfig{
			NetworksAllowed:          e.NetworksAllowed,
			NetworksAllowAll:         e.NetworksAllowAll,
			NetworkAddressesAllowed:  e.NetworkAddressesAllowed,
			NetworkAddressesAllowAll: e.NetworkAddressesAllowAll,
		}))
		functions = append(functions, wnet.ConnRead())
		functions = append(functions, wnet.ConnWrite())
		functions = append(functions, wnet.ConnClose())
	}

	functions = append(functions, e.HostFunctions...)
	return functions
}
