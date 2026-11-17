package market

type Asset struct {
	ID            int64  `json:"id"`
	Symbol        string `json:"symbol"`
	Name          string `json:"name"`
	DecimalPlaces int32  `json:"decimal_places"`
	IsActive      bool   `json:"is_active"`
}

type Market struct {
	ID                int64  `json:"id"`
	Symbol            string `json:"symbol"`
	BaseAsset         string `json:"base_asset"`
	QuoteAsset        string `json:"quote_asset"`
	PricePrecision    int32  `json:"price_precision"`
	QuantityPrecision int32  `json:"quantity_precision"`
	ReferencePrice    string `json:"reference_price"`
	PriceSource       string `json:"price_source"`
	Status            string `json:"status"`
}
