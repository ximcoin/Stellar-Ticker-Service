package scraper

//func Test_ScraperConfig_FetchAllTrades_doesntCrashWhenReceivesAnError(t *testing.T) {
//	horizonClient := &horizonclient.MockClient{}
//	horizonClient.
//		On("Trades", horizonclient.TradeRequest{Limit: 200, Order: horizonclient.OrderDesc}).
//		Return(hProtocol.TradesPage{}, errors.New("something went wrong"))
//
//	sc := ScraperConfig{
//		Logger: log.DefaultLogger,
//		Client: horizonClient,
//	}
//
//	trades, err := sc.FetchAllTrades(time.Now(), 0)
//	assert.EqualError(t, err, "something went wrong")
//	assert.Empty(t, trades)
//}
//
//func Test_ScraperConfig_FetchAllTrades_doesntCrashWhenReceivesEmptyList(t *testing.T) {
//	horizonClient := &horizonclient.MockClient{}
//	horizonClient.
//		On("Trades", horizonclient.TradeRequest{Limit: 200, Order: horizonclient.OrderDesc}).
//		Return(hProtocol.TradesPage{}, nil)
//
//	sc := ScraperConfig{
//		Logger: log.DefaultLogger,
//		Client: horizonClient,
//	}
//
//	trades, err := sc.FetchAllTrades(time.Now(), 0)
//	assert.NoError(t, err)
//	assert.Empty(t, trades)
//}
