package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"slices"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"github.com/mymmrac/wape"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:                   "wape [-e env.toml] [-f func] [-i input] [file.wasm]",
	Short:                 "Run WASM modules with WAPE environment",
	Args:                  cobra.MaximumNArgs(1),
	Version:               version(),
	RunE:                  run,
	SilenceErrors:         true,
	DisableFlagsInUseLine: true,
}

var (
	envFilepath  string
	funcName     string
	input        string
	debugEnabled bool

	randSourceFromHost bool
)

func init() {
	rootCmd.Flags().StringVarP(&envFilepath, "env", "e", "", "WAPE environment file (toml)")
	rootCmd.Flags().StringVarP(&funcName, "func", "f", "main", "function to call")
	rootCmd.Flags().StringVarP(&input, "input", "i", "", "input data")
	rootCmd.Flags().BoolVar(&debugEnabled, "debug", false, "debug mode")
	rootCmd.Flags().BoolVar(&randSourceFromHost, "rand-host", false, "random from host")
	_ = rootCmd.Flags().MarkHidden("debug")
}

func run(cmd *cobra.Command, args []string) error {
	var env *wape.Environment

	if envFilepath != "" {
		envFile, err := os.ReadFile(envFilepath)
		if err != nil {
			return fmt.Errorf("failed to read the environment file: %w", err)
		}

		if err = toml.Unmarshal(envFile, &env); err != nil {
			return fmt.Errorf("failed to parse the environment file: %w", err)
		}
	} else {
		env = wape.NewEnvironment()
	}

	if len(args) > 0 {
		if slices.ContainsFunc(env.Modules, func(data wape.ModuleData) bool {
			return data.Name == "main"
		}) {
			return fmt.Errorf("main module already provided in the environment")
		}

		env.Modules = append(env.Modules, wape.ModuleData{
			Name: "main",
			File: args[0],
		})
	}

	if len(env.Modules) == 0 {
		return fmt.Errorf("no WASM modules specified")
	}

	if randSourceFromHost {
		env.RandSourceFromHost = true
	}

	if debugEnabled {
		envData, err := json.MarshalIndent(env, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal the environment: %w", err)
		}
		fmt.Println("Environment:")
		fmt.Println(string(envData))
		fmt.Println()
	}

	ctx := cmd.Context()

	plugin, err := wape.NewPlugin(ctx, env)
	if err != nil {
		return fmt.Errorf("failed to create the plugin: %w", err)
	}

	exit, output, err := plugin.CallWithContext(ctx, funcName, []byte(input))
	if err != nil {
		return fmt.Errorf("failed to call the function: %w", err)
	}

	if len(output) > 0 {
		fmt.Println(string(output))
	}

	if exit != 0 {
		return fmt.Errorf("function returned non-zero exit code: %d", exit)
	}

	return nil
}

func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	return fmt.Sprintf("%s build with %s", info.Main.Version, info.GoVersion)
}
