package ticker

import (
	"github.com/stellar/go/services/ticker/internal/scraper"
)

// MarketSummary represents a summary of statistics of all valid markets
// within a given period of time.
type MarketSummary struct {
	GeneratedAt        int64         `json:"generated_at"`
	GeneratedAtRFC3339 string        `json:"generated_at_rfc3339"`
	Pairs              []MarketStats `json:"pairs"`
}

// MarketStats represents the statistics of a specific market (identified by
// a trade pair).
type MarketStats struct {
	TradePairName    string  `json:"name"`
	BaseVolume24h    float64 `json:"base_volume"`
	CounterVolume24h float64 `json:"counter_volume"`
	TradeCount24h    int64   `json:"trade_count"`
	Open24h          float64 `json:"open"`
	Low24h           float64 `json:"low"`
	High24h          float64 `json:"high"`
	Change24h        float64 `json:"change"`
	BaseVolume7d     float64 `json:"base_volume_7d"`
	CounterVolume7d  float64 `json:"counter_volume_7d"`
	TradeCount7d     int64   `json:"trade_count_7d"`
	Open7d           float64 `json:"open_7d"`
	Low7d            float64 `json:"low_7d"`
	High7d           float64 `json:"high_7d"`
	Change7d         float64 `json:"change_7d"`
	Price            float64 `json:"price"`
	Close            float64 `json:"close"`
	CloseTime        string  `json:"close_time"`
	BidCount         int     `json:"bid_count"`
	BidVolume        float64 `json:"bid_volume"`
	BidMax           float64 `json:"bid_max"`
	AskCount         int     `json:"ask_count"`
	AskVolume        float64 `json:"ask_volume"`
	AskMin           float64 `json:"ask_min"`
	Spread           float64 `json:"spread"`
	SpreadMidPoint   float64 `json:"spread_mid_point"`
}

// PartialMarketSummary represents a summary of statistics of all valid markets
// within a given period of time.
type PartialMarketSummary struct {
	GeneratedAt        int64                `json:"generated_at"`
	GeneratedAtRFC3339 string               `json:"generated_at_rfc3339"`
	Pairs              []PartialMarketStats `json:"pairs"`
}

// PartialMarketStats represents the statistics of a specific market (identified by
// a trade pair).
type PartialMarketStats struct {
	TradePairName      string  `json:"name"`
	BaseAssetID        int32   `json:"base_asset_id"`
	BaseAssetCode      string  `json:"base_asset_code"`
	BaseAssetIssuer    string  `json:"base_asset_issuer"`
	BaseAssetType      string  `json:"base_asset_type"`
	CounterAssetID     int32   `json:"counter_asset_id"`
	CounterAssetCode   string  `json:"counter_asset_code"`
	CounterAssetIssuer string  `json:"counter_asset_issuer"`
	CounterAssetType   string  `json:"counter_asset_type"`
	BaseVolume         float64 `json:"base_volume"`
	CounterVolume      float64 `json:"counter_volume"`
	TradeCount         int32   `json:"trade_count"`
	Open               float64 `json:"open"`
	Low                float64 `json:"low"`
	High               float64 `json:"high"`
	Change             float64 `json:"change"`
	Close              float64 `json:"close"`
	NumBids            int     `json:"num_bids"`
	BidVolume          float64 `json:"bid_volume"`
	HighestBid         float64 `json:"highest_bid"`
	NumAsks            int     `json:"num_asks"`
	AskVolume          float64 `json:"ask_volume"`
	LowestAsk          float64 `json:"lowest_ask"`
	Spread             float64 `json:"spread"`
	SpreadMidPoint     float64 `json:"spread_mid_point"`
}

// Asset Summary represents the collection of valid assets.
type AssetSummary struct {
	GeneratedAt        int64   `json:"generated_at"`
	GeneratedAtRFC3339 string  `json:"generated_at_rfc3339"`
	Assets             []Asset `json:"assets"`
}

// Asset represent the aggregated data for a given asset.
type Asset struct {
	scraper.FinalAsset

	IssuerDetail       Issuer `json:"issuer_detail"`
	LastValidTimestamp string `json:"last_valid"`
}

// Issuer represents the aggregated data for a given issuer.
type Issuer struct {
	PublicKey        string `json:"public_key"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	TOMLURL          string `json:"toml_url"`
	FederationServer string `json:"federation_server"`
	AuthServer       string `json:"auth_server"`
	TransferServer   string `json:"transfer_server"`
	WebAuthEndpoint  string `json:"web_auth_endpoint"`
	DepositServer    string `json:"deposit_server"`
	OrgTwitter       string `json:"org_twitter"`
}
