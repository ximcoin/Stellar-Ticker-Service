package scraper

import (
	"context"
	"errors"
	"fmt"
	hProtocol "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/services/ticker/internal/tickerdb"
	hlog "github.com/stellar/go/support/log"
	"math/big"
	"strconv"
)

// TODO: 30 sec for an insert of 100k records -> 12k rows?
func PersistTrades(
	ctx context.Context,
	s *tickerdb.TickerSession,
	l *hlog.Entry,
	trades []hProtocol.Trade,
) error {
	var dbTrades []tickerdb.Trade

	for _, trade := range trades {
		var bID, cID int32
		bID, cID, err := FindBaseAndCounter(ctx, s, trade)
		if err != nil {
			continue
		}

		var dbTrade tickerdb.Trade
		dbTrade, err = HProtocolTradeToDBTrade(trade, bID, cID)
		if err != nil {
			l.Error("Could not convert entry to DB Trade: ", err)
			continue
		}
		dbTrades = append(dbTrades, dbTrade)
	}

	l.Infof("Inserting %d entries in the database.\n", len(dbTrades))
	if err := s.BulkInsertTrades(ctx, dbTrades); err != nil {
		fmt.Println(err)
	}

	return nil
}

// FindBaseAndCounter tries to find the Base and Counter assets IDs in the database,
// and returns an error if it doesn't find any.
func FindBaseAndCounter(ctx context.Context, s *tickerdb.TickerSession, trade hProtocol.Trade) (bID int32, cID int32, err error) {
	bFound, bID, err := s.GetAssetByCodeAndIssuerAccount(
		ctx,
		trade.BaseAssetCode,
		trade.BaseAssetIssuer,
	)
	if err != nil {
		return
	}

	cFound, cID, err := s.GetAssetByCodeAndIssuerAccount(
		ctx,
		trade.CounterAssetCode,
		trade.CounterAssetIssuer,
	)
	if err != nil {
		return
	}

	if !bFound || !cFound {
		err = errors.New("base or counter asset no found")
		return
	}

	return
}

// HProtocolTradeToDBTrade converts from a hProtocol.Trade to a tickerdb.Trade
func HProtocolTradeToDBTrade(
	hpt hProtocol.Trade,
	baseAssetID int32,
	counterAssetID int32,
) (trade tickerdb.Trade, err error) {
	fBaseAmount, err := strconv.ParseFloat(hpt.BaseAmount, 64)
	if err != nil {
		return
	}
	fCounterAmount, err := strconv.ParseFloat(hpt.CounterAmount, 64)
	if err != nil {
		return
	}

	rPrice := big.NewRat(hpt.Price.D, hpt.Price.N)
	fPrice, _ := rPrice.Float64()

	trade = tickerdb.Trade{
		HorizonID:       hpt.ID,
		LedgerCloseTime: hpt.LedgerCloseTime,
		OfferID:         hpt.OfferID,
		BaseOfferID:     hpt.BaseOfferID,
		BaseAccount:     hpt.BaseAccount,
		BaseAmount:      fBaseAmount,
		BaseAssetID:     baseAssetID,
		CounterOfferID:  hpt.CounterOfferID,
		CounterAccount:  hpt.CounterAccount,
		CounterAmount:   fCounterAmount,
		CounterAssetID:  counterAssetID,
		BaseIsSeller:    hpt.BaseIsSeller,
		Price:           fPrice,
	}

	return
}
