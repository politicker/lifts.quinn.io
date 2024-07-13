package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/politicker/lifts.quinn.io/internal/cmdutil"
	"github.com/politicker/lifts.quinn.io/internal/db"
	"github.com/politicker/lifts.quinn.io/internal/domain"
	"github.com/spf13/cobra"
)

func ExtractCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract",
		Args:  cobra.ExactArgs(0),
		Short: "Import lift history from a CSV file",
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath, err := cmd.Flags().GetString("file")
			if err != nil {
				return err
			}

			if filePath == "" {
				return errors.New("file is a required flag")
			}

			csvFile, err := os.ReadFile(filePath)
			if err != nil {
				panic(err)
			}

			if len(csvFile) == 0 {
				return errors.New("lifts CSV file is empty")
			}

			logger := cmdutil.NewLogger("import")
			defer func() { _ = logger.Sync() }()

			database, err := cmdutil.NewDBConnection(ctx)
			if err != nil {
				return err
			}
			defer database.Close()

			queries := db.New(database)

			importer := domain.NewImporter(logger, queries)
			err = importer.Run(ctx, csvFile)
			if err != nil {
				panic(err)
			}

			logger.Info("imported workout history")
			return nil
		},
	}

	cmd.Flags().StringP("file", "f", "", "The CSV file to import")
	return cmd
}
