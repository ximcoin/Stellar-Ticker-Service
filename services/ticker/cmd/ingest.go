package cmd

import (
	"bufio"
	"context"
	"os"

	"github.com/lib/pq"
	"github.com/spf13/cobra"
	ticker "github.com/stellar/go/services/ticker/internal"
	"github.com/stellar/go/services/ticker/internal/tickerdb"
)

var ShouldStream bool
var BackfillHours int

func init() {
	rootCmd.AddCommand(cmdIngest)
	//cmdIngest.AddCommand(cmdIngestAssets)
	cmdIngest.AddCommand(cmdIngestFilteredAssets)
	cmdIngest.AddCommand(cmdIngestTrades)
	cmdIngest.AddCommand(cmdIngestFilteredTrades)
	//cmdIngest.AddCommand(cmdIngestOrderbooks)
	cmdIngest.AddCommand(cmdIngestFilteredOrderbooks)

	cmdIngestTrades.Flags().BoolVar(
		&ShouldStream,
		"stream",
		false,
		"Continuously stream new trades from the Horizon Stream API as a daemon",
	)

	cmdIngestTrades.Flags().IntVar(
		&BackfillHours,
		"num-hours",
		1*24, //with 24h, uses roughly 3GB of memory.
		"Number of past hours to backfill trade data",
	)

	cmdIngestFilteredAssets.Flags().StringVarP(
		&filePath,
		"file",
		"f",
		"",
		"Filter assets by issuers defined in a file",
	)

	cmdIngestFilteredTrades.Flags().StringVarP(
		&filePath,
		"file",
		"f",
		"",
		"Filter trades by issuers defined in a file",
	)

	cmdIngestFilteredOrderbooks.Flags().StringVarP(
		&filePath,
		"file",
		"f",
		"",
		"Filter orderbooks by issuers defined in a file",
	)
}

var cmdIngest = &cobra.Command{
	Use:   "ingest [data type]",
	Short: "Ingests new data from data type into the database.",
}

var cmdIngestAssets = &cobra.Command{
	Use:   "assets",
	Short: "Refreshes the asset database with new data retrieved from Horizon.",
	Run: func(cmd *cobra.Command, args []string) {
		Logger.Info("Refreshing the asset database")
		dbInfo, err := pq.ParseURL(DatabaseURL)
		if err != nil {
			Logger.Fatal("could not parse db-url:", err)
		}

		session, err := tickerdb.CreateSession("postgres", dbInfo)
		if err != nil {
			Logger.Fatal("could not connect to db:", err)
		}
		defer session.DB.Close()

		ctx := context.Background()
		err = ticker.RefreshAssets(ctx, &session, Client, Logger)
		if err != nil {
			Logger.Fatal("could not refresh asset database:", err)
		}
	},
}

var cmdIngestFilteredAssets = &cobra.Command{
	Use:   "filtered-assets",
	Short: "Refreshes the asset database with new data retrieved from Horizon filtered by issuers defined in a file.",
	Run: func(cmd *cobra.Command, args []string) {
		if filePath == "" {
			Logger.Fatal("file flag is required")
		}

		Logger.Info("Refreshing the asset database")
		dbInfo, err := pq.ParseURL(DatabaseURL)
		if err != nil {
			Logger.Fatal("could not parse db-url:", err)
		}

		session, err := tickerdb.CreateSession("postgres", dbInfo)
		if err != nil {
			Logger.Fatal("could not connect to db:", err)
		}
		defer session.DB.Close()

		fileContents, err := getIssuers(filePath)
		if err != nil {
			Logger.Fatal("could not get issuers from file:", err)
		}

		// deduplicate the file contents
		issuers := removeDuplicate(fileContents)

		ctx := context.Background()

		// for each issuer we spawn 20 threads

		// loop over the issuers and refresh the assets
		for _, issuer := range issuers {
			Logger.Infof("Refreshing assets for issuer: %s", issuer)
			err = ticker.RefreshFilteredAssets(ctx, &session, Client, Logger, issuer)
			if err != nil {
				Logger.Fatal("could not refresh asset database:", err)
			}
		}
	},
}

var cmdIngestTrades = &cobra.Command{
	Use:   "trades",
	Short: "Fills the trade database with data retrieved from Horizon.",
	Run: func(cmd *cobra.Command, args []string) {
		dbInfo, err := pq.ParseURL(DatabaseURL)
		if err != nil {
			Logger.Fatal("could not parse db-url:", err)
		}

		session, err := tickerdb.CreateSession("postgres", dbInfo)
		if err != nil {
			Logger.Fatal("could not connect to db:", err)
		}
		defer session.DB.Close()

		ctx := context.Background()
		numDays := float32(BackfillHours) / 24.0
		Logger.Infof(
			"Backfilling Trade data for the past %d hour(s) [%.2f days]\n",
			BackfillHours,
			numDays,
		)
		err = ticker.BackfillTrades(ctx, &session, Client, Logger, BackfillHours, 0)
		if err != nil {
			Logger.Fatal("could not refresh trade database:", err)
		}

		if ShouldStream {
			Logger.Info("Streaming new data (this is a continuous process)")
			err = ticker.StreamTrades(ctx, &session, Client, Logger)
			if err != nil {
				Logger.Fatal("could not refresh trade database:", err)
			}
		}
	},
}

var cmdIngestFilteredTrades = &cobra.Command{
	Use:   "filtered-trades",
	Short: "Fills the trade database with data retrieved from Horizon filtered by issuers defined in a file.",
	Run: func(cmd *cobra.Command, args []string) {
		if filePath == "" {
			Logger.Fatal("file flag is required")
		}

		dbInfo, err := pq.ParseURL(DatabaseURL)
		if err != nil {
			Logger.Fatal("could not parse db-url:", err)
		}

		session, err := tickerdb.CreateSession("postgres", dbInfo)
		if err != nil {
			Logger.Fatal("could not connect to db:", err)
		}
		defer session.DB.Close()

		ctx := context.Background()
		numDays := float32(BackfillHours) / 24.0
		Logger.Infof(
			"Backfilling Trade data for the past %d hour(s) [%.2f days]\n",
			BackfillHours,
			numDays,
		)

		fileContents, err := getIssuers(filePath)
		if err != nil {
			Logger.Fatal("could not get issuers from file:", err)
		}

		// deduplicate the file contents
		issuers := removeDuplicate(fileContents)

		// loop over the issuers and refresh the trades
		for _, issuer := range issuers {
			err = ticker.BackfillFilteredTrades(ctx, &session, Client, Logger, BackfillHours, 0, issuer)
			if err != nil {
				Logger.Fatal("could not refresh trade database:", err)
			}
		}

		if ShouldStream {
			Logger.Info("Streaming new data (this is a continuous process)")
			err = ticker.StreamTrades(ctx, &session, Client, Logger)
			if err != nil {
				Logger.Fatal("could not refresh trade database:", err)
			}
		}
	},
}

var cmdIngestOrderbooks = &cobra.Command{
	Use:   "orderbooks",
	Short: "Refreshes the orderbook stats database with new data retrieved from Horizon.",
	Run: func(cmd *cobra.Command, args []string) {
		Logger.Info("Refreshing the asset database")
		dbInfo, err := pq.ParseURL(DatabaseURL)
		if err != nil {
			Logger.Fatal("could not parse db-url:", err)
		}

		session, err := tickerdb.CreateSession("postgres", dbInfo)
		if err != nil {
			Logger.Fatal("could not connect to db:", err)
		}
		defer session.DB.Close()

		err = ticker.RefreshOrderbookEntries(&session, Client, Logger)
		if err != nil {
			Logger.Fatal("could not refresh orderbook database:", err)
		}
	},
}

var cmdIngestFilteredOrderbooks = &cobra.Command{
	Use:   "filtered-orderbooks",
	Short: "Refreshes the orderbook stats database with new data retrieved from Horizon filtered by issuers defined in a file.",
	Run: func(cmd *cobra.Command, args []string) {
		if filePath == "" {
			Logger.Fatal("file flag is required")
		}

		Logger.Info("Refreshing the asset database")
		dbInfo, err := pq.ParseURL(DatabaseURL)
		if err != nil {
			Logger.Fatal("could not parse db-url:", err)
		}

		session, err := tickerdb.CreateSession("postgres", dbInfo)
		if err != nil {
			Logger.Fatal("could not connect to db:", err)
		}
		defer session.DB.Close()

		fileContents, err := getIssuers(filePath)
		if err != nil {
			Logger.Fatal("could not get issuers from file:", err)
		}

		// deduplicate the file contents
		issuers := removeDuplicate(fileContents)

		err = ticker.RefreshFilteredOrderbookEntries(&session, Client, Logger, issuers)
		if err != nil {
			Logger.Fatal("could not refresh orderbook database:", err)
		}
	},
}

func getIssuers(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var issuers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		issuers = append(issuers, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return issuers, nil
}

func removeDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
