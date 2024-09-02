package scraper

import (
	"fmt"
	"github.com/stellar/go/clients/horizonclient"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	hProtocol "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/log"
	hlog "github.com/stellar/go/support/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This test showed that EURC is going to the db
func TestAssetScraper_EURC(t *testing.T) {
	sc := ScraperConfig{
		Client: horizonclient.DefaultPublicNetClient,
		Logger: hlog.New(),
	}
	assets, err := sc.retrieveFilteredAssets(0, "GDHU6WRG4IEQXM5NZ4BMPKOXHW76MZM4Y2IEMFDVXBSDP6SJY4ITNPP2")
	t.Log(assets, err)

	//assetQueue := make(chan any, 20)
	// we get the asset successfully.

	// [{{{https://circle.com/.well-known/stellar.toml false}} {credit_alphanum4 EURC GDHU6WRG4IEQXM5NZ4BMPKOXHW76MZM4Y2IEMFDVXBSDP6SJY4ITNPP2} EURC_GDHU6WRG4IEQXM5NZ4BMPKOXHW76MZM4Y2IEMFDVXBSDP6SJY4ITNPP2_credit_alphanum4 CDTKPWPLOURQA2SGTKTUQOWRCBZEORB4BWBOMJ3D3ZTQQSGE5F6JBQLV 6592 60 41 10 1 1577540.1109837 {6592 0 0} 1279.3326979 97607.4874612 76556.4923921 1.0913390 {1577540.1109837 0.0000000 0.0000000} {false true false false}}]

	//nonTrash, trash := sc.parallelProcessAssets(assets, 20, assetQueue)
	//t.Log(nonTrash, trash)
}

func TestShouldDiscardAsset(t *testing.T) {
	testAsset := hProtocol.AssetStat{
		Amount: "",
	}

	assert.Equal(t, shouldDiscardAsset(testAsset, true), true)

	testAsset = hProtocol.AssetStat{
		Amount: "0.0",
	}
	assert.Equal(t, shouldDiscardAsset(testAsset, true), true)

	testAsset = hProtocol.AssetStat{
		Amount: "0",
	}
	assert.Equal(t, shouldDiscardAsset(testAsset, true), true)

	testAsset = hProtocol.AssetStat{
		Amount:      "123901.0129310",
		NumAccounts: 8,
	}
	assert.Equal(t, shouldDiscardAsset(testAsset, true), true)

	testAsset = hProtocol.AssetStat{
		Amount:      "123901.0129310",
		NumAccounts: 12,
	}
	testAsset.Code = "REMOVE"
	assert.Equal(t, shouldDiscardAsset(testAsset, true), true)

	testAsset = hProtocol.AssetStat{
		Amount:      "123901.0129310",
		NumAccounts: 100,
	}
	testAsset.Code = "SOMETHINGVALID"
	testAsset.Links.Toml.Href = ""
	assert.Equal(t, shouldDiscardAsset(testAsset, true), false)

	testAsset = hProtocol.AssetStat{
		Amount:      "123901.0129310",
		NumAccounts: 40,
	}
	testAsset.Code = "SOMETHINGVALID"
	testAsset.Links.Toml.Href = "http://www.stellar.org/.well-known/stellar.toml"
	assert.Equal(t, shouldDiscardAsset(testAsset, true), true)

	testAsset = hProtocol.AssetStat{
		Amount:      "123901.0129310",
		NumAccounts: 40,
	}
	testAsset.Code = "SOMETHINGVALID"
	testAsset.Links.Toml.Href = ""
	assert.Equal(t, shouldDiscardAsset(testAsset, true), true)

	testAsset = hProtocol.AssetStat{
		Amount:      "123901.0129310",
		NumAccounts: 40,
	}
	testAsset.Code = "SOMETHINGVALID"
	testAsset.Links.Toml.Href = "https://www.stellar.org/.well-known/stellar.toml"
	assert.Equal(t, shouldDiscardAsset(testAsset, true), false)
}

func TestDomainsMatch(t *testing.T) {
	tomlURL, _ := url.Parse("https://stellar.org/stellar.toml")
	orgURL, _ := url.Parse("https://stellar.org/")
	assert.True(t, domainsMatch(tomlURL, orgURL))

	tomlURL, _ = url.Parse("https://assets.stellar.org/stellar.toml")
	orgURL, _ = url.Parse("https://stellar.org/")
	assert.False(t, domainsMatch(tomlURL, orgURL))

	tomlURL, _ = url.Parse("https://stellar.org/stellar.toml")
	orgURL, _ = url.Parse("https://home.stellar.org/")
	assert.True(t, domainsMatch(tomlURL, orgURL))

	tomlURL, _ = url.Parse("https://stellar.org/stellar.toml")
	orgURL, _ = url.Parse("https://home.stellar.com/")
	assert.False(t, domainsMatch(tomlURL, orgURL))

	tomlURL, _ = url.Parse("https://stellar.org/stellar.toml")
	orgURL, _ = url.Parse("https://stellar.com/")
	assert.False(t, domainsMatch(tomlURL, orgURL))
}

func TestIsDomainVerified(t *testing.T) {
	tomlURL := "https://stellar.org/stellar.toml"
	orgURL := "https://stellar.org/"
	hasCurrency := true
	assert.True(t, isDomainVerified(orgURL, tomlURL, hasCurrency))

	tomlURL = "https://stellar.org/stellar.toml"
	orgURL = ""
	hasCurrency = true
	assert.True(t, isDomainVerified(orgURL, tomlURL, hasCurrency))

	tomlURL = ""
	orgURL = ""
	hasCurrency = true
	assert.False(t, isDomainVerified(orgURL, tomlURL, hasCurrency))

	tomlURL = "https://stellar.org/stellar.toml"
	orgURL = "https://stellar.org/"
	hasCurrency = false
	assert.False(t, isDomainVerified(orgURL, tomlURL, hasCurrency))

	tomlURL = "http://stellar.org/stellar.toml"
	orgURL = "https://stellar.org/"
	hasCurrency = true
	assert.False(t, isDomainVerified(orgURL, tomlURL, hasCurrency))

	tomlURL = "https://stellar.org/stellar.toml"
	orgURL = "http://stellar.org/"
	hasCurrency = true
	assert.False(t, isDomainVerified(orgURL, tomlURL, hasCurrency))

	tomlURL = "https://stellar.org/stellar.toml"
	orgURL = "https://stellar.com/"
	hasCurrency = true
	assert.False(t, isDomainVerified(orgURL, tomlURL, hasCurrency))
}

func TestIgnoreInvalidTOMLUrls(t *testing.T) {
	invalidURL := "https:// there is something wrong here.com/stellar.toml"
	_, err := fetchTOMLData(invalidURL)

	urlErr, ok := errors.Cause(err).(*url.Error)
	if !ok {
		t.Fatalf("err expected to be a url.Error but was %#v", err)
	}
	assert.Equal(t, "parse", urlErr.Op)
	assert.Equal(t, "https:// there is something wrong here.com/stellar.toml", urlErr.URL)
	assert.EqualError(t, urlErr.Err, `invalid character " " in host name`)
}

func TestProcessAsset_notCached(t *testing.T) {
	logger := log.DefaultLogger
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `SIGNING_KEY="not cached signing key"`)
	}))
	asset := hProtocol.AssetStat{
		Amount:      "123901.0129310",
		NumAccounts: 100,
	}
	asset.Code = "SOMETHINGVALID"
	asset.Links.Toml.Href = server.URL
	tomlCache := &TOMLCache{}
	finalAsset, err := processAsset(logger, asset, tomlCache, true)
	require.NoError(t, err)
	assert.NotZero(t, finalAsset)
	assert.Equal(t, "not cached signing key", finalAsset.IssuerDetails.SigningKey)
	cachedTOML, ok := tomlCache.Get(server.URL)
	assert.True(t, ok)
	assert.Equal(t, TOMLIssuer{SigningKey: "not cached signing key"}, cachedTOML)
}

func TestProcessAsset_cached(t *testing.T) {
	logger := log.DefaultLogger
	asset := hProtocol.AssetStat{
		Amount:      "123901.0129310",
		NumAccounts: 100,
	}
	asset.Code = "SOMETHINGVALID"
	asset.Links.Toml.Href = "url"
	tomlCache := &TOMLCache{}
	tomlCache.Set("url", TOMLIssuer{SigningKey: "signing key"})
	finalAsset, err := processAsset(logger, asset, tomlCache, true)
	require.NoError(t, err)
	assert.NotZero(t, finalAsset)
	assert.Equal(t, "signing key", finalAsset.IssuerDetails.SigningKey)
}
