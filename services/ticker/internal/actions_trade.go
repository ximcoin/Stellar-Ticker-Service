package ticker

import (
	"context"
	"fmt"
	"time"

	horizonclient "github.com/stellar/go/clients/horizonclient"
	hProtocol "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/services/ticker/internal/scraper"
	"github.com/stellar/go/services/ticker/internal/tickerdb"
	hlog "github.com/stellar/go/support/log"
)

// StreamTrades constantly streams and ingests new trades directly from horizon.
func StreamTrades(
	ctx context.Context,
	s *tickerdb.TickerSession,
	c *horizonclient.Client,
	l *hlog.Entry,
) error {
	sc := scraper.ScraperConfig{
		Client: c,
		Logger: l,
		Ctx:    &ctx,
	}
	handler := func(trade hProtocol.Trade) {
		l.Infof("New trade arrived. ID: %v; Close Time: %v\n", trade.ID, trade.LedgerCloseTime)
		scraper.NormalizeTradeAssets(&trade)
		bID, cID, err := scraper.FindBaseAndCounter(ctx, s, trade)
		if err != nil {
			l.Error(err)
			return
		}
		dbTrade, err := scraper.HProtocolTradeToDBTrade(trade, bID, cID)
		if err != nil {
			l.Error(err)
			return
		}

		err = s.BulkInsertTrades(ctx, []tickerdb.Trade{dbTrade})
		if err != nil {
			l.Error("Could not insert trade in database: ", trade.ID)
		}
	}

	// Ensure we start streaming from the last stored trade
	lastTrade, err := s.GetLastTrade(ctx)
	if err != nil {
		return err
	}

	cursor := lastTrade.HorizonID
	return sc.StreamNewTrades(cursor, handler)
}

// BackfillTrades ingest the most recent trades (limited to numDays) directly from Horizon
// into the database.
func BackfillTrades(
	ctx context.Context,
	s *tickerdb.TickerSession,
	c *horizonclient.Client,
	l *hlog.Entry,
	numHours int,
	limit int,
) error {
	sc := scraper.ScraperConfig{
		Client: c,
		Logger: l,
	}
	now := time.Now()
	since := now.Add(time.Hour * -time.Duration(numHours))
	trades, err := sc.FetchAllTrades(ctx, s, l, since, limit)
	if err != nil {
		return err
	}

	return scraper.PersistTrades(ctx, s, l, trades)
}

// BackfillFilteredTrades ingest the most recent trades (limited to numDays) directly from Horizon
// into the database, filtered by issuer.
func BackfillFilteredTrades(
	ctx context.Context,
	s *tickerdb.TickerSession,
	c *horizonclient.Client,
	l *hlog.Entry,
	numHours int,
	limit int,
	issuer string,
) error {
	sc := scraper.ScraperConfig{
		Client: c,
		Logger: l,
	}
	now := time.Now()
	since := now.Add(time.Hour * -time.Duration(numHours))
	trades, err := sc.FetchFilteredTrades(since, limit, issuer)
	if err != nil {
		return err
	}

	var dbTrades []tickerdb.Trade

	for _, trade := range trades {
		var bID, cID int32
		bID, cID, err = scraper.FindBaseAndCounter(ctx, s, trade)
		if err != nil {
			continue
		}

		var dbTrade tickerdb.Trade
		dbTrade, err = scraper.HProtocolTradeToDBTrade(trade, bID, cID)
		if err != nil {
			l.Error("Could not convert entry to DB Trade: ", err)
			continue
		}
		dbTrades = append(dbTrades, dbTrade)
	}

	l.Infof("Inserting %d entries in the database.\n", len(dbTrades))
	err = s.BulkInsertTrades(ctx, dbTrades)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
