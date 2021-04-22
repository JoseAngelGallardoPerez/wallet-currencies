package ratesupdater

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/Confialink/wallet-currencies/internal/services/exchange"
)

const (
	xmlRatesPath = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"
	regexPattern = `<Cube currency='([A-Z]{3})' rate='(\d+(\.\d+)?)'\/>`
)

var (
	compiledPattern = regexp.MustCompile(regexPattern)
)

type Ecb struct {
}

func NewEcb() CurrencySourceFactory {
	return &Ecb{}
}

// Init uploads rates from ECB and prepare them
func (s *Ecb) Init() (exchange.RateSource, error) {
	report, err := s.downloadReport()
	if err != nil {
		return nil, errors.Wrap(err, "cannot download ECB rates")
	}

	directSource := exchange.NewDirectRateSource()

	for currencyCode, value := range s.parseRates(report) {
		decimalVal, err := decimal.NewFromString(value)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot convert to decimal value %s for currency %s", value, currencyCode)
		}
		_ = directSource.Set(exchange.NewRate("EUR", currencyCode, decimalVal))
	}

	return exchange.NewPivotRateSource("EUR", directSource), nil
}

// downloadReport downloads currency rates
func (s *Ecb) downloadReport() (xmlString string, err error) {
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Get(xmlRatesPath)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	xmlBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return string(xmlBytes), nil
}

// parseRates return map like {"USD": "1.111"}
func (s *Ecb) parseRates(report string) map[string]string {
	rates := make(map[string]string)
	matches := compiledPattern.FindAllStringSubmatch(report, -1)
	for _, match := range matches {
		rates[match[1]] = match[2]
	}
	return rates
}
