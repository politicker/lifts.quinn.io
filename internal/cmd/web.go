package cmd

import (
	"context"
	"os"
	"strconv"

	"github.com/politicker/lifts.quinn.io/internal/cmdutil"
	"github.com/politicker/lifts.quinn.io/internal/web"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func WebCmd(ctx context.Context) *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "web",
		Args:  cobra.ExactArgs(0),
		Short: "Start the web server",
		RunE: func(cmd *cobra.Command, args []string) error {
			port = 8001
			if os.Getenv("PORT") != "" {
				port, _ = strconv.Atoi(os.Getenv("PORT"))
			}

			logger := cmdutil.NewLogger("web")
			defer func() { _ = logger.Sync() }()

			db, err := cmdutil.NewDBConnection(ctx)
			if err != nil {
				return err
			}
			defer db.Close()

			srv := web.NewWeb(ctx, logger, db, port)
			srv.Start()

			logger.Info("web server started", zap.Int("port", port))

			<-ctx.Done()
			return nil
		},
	}

	return cmd
}
