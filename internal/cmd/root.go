package cmd

import (
	"context"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/spf13/cobra"
)

func Execute(ctx context.Context) int {

	profile := false

	rootCmd := &cobra.Command{
		Use:   "lifts.quinn.io",
		Short: "A web service for keeping up with Quinn's lifts",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !profile {
				return nil
			}

			f, perr := os.Create("cpu.pprof")
			if perr != nil {
				return perr
			}

			_ = pprof.StartCPUProfile(f)
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if !profile {
				return nil
			}

			pprof.StopCPUProfile()

			f, perr := os.Create("mem.pprof")
			if perr != nil {
				return perr
			}
			defer f.Close()

			runtime.GC()
			err := pprof.WriteHeapProfile(f)
			return err
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&profile, "profile", "p", false, "record CPU pprof")

	rootCmd.AddCommand(WebCmd(ctx))
	rootCmd.AddCommand(ExtractCmd(ctx))

	// I'm not sure what this is for.
	// go func() {
	// 	_ = http.ListenAndServe("localhost:6060", nil)
	// }()

	if err := rootCmd.Execute(); err != nil {
		return 1
	}

	return 0
}
