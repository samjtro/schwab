package trader

type Transaction struct {
	ActivityId     int
	Time           string
	User           User
	Description    string
	AccountNumber  string
	Type           string
	Status         string
	SubAccount     string
	TradeDate      string
	SettlementDate string
	PositionId     int
	OrderId        int
	NetAmount      int
	ActivityType   string
	TransferItems  TransferItems
}

type User struct {
	CdDomainId     string
	Login          string
	Type           string
	UserId         int
	SystemUserName string
	FirstName      string
	LastName       string
	BrokerRepCode  string
}

type TransferItems struct {
	Instrument     InstrumentRef
	Amount         int
	Cost           int
	Price          int
	FeeType        string
	PositionEffect string
}

type InstrumentRef struct {
	Cusip        string
	Symbol       string
	Description  string
	InstrumentId int
	NetChange    int
	Type         string
}

type AccountNumbers struct {
	AccountNumber string
	HashValue     string
}

type Account struct {
	SecuritiesAccount SecuritiesAccount
	AggregatedBalance AggregatedBalance
}

type SecuritiesAccount struct {
	Type                    string
	AccountNumber           string
	RoundTrips              int
	IsDayTrader             bool
	IsClosingOnlyRestricted bool
	PFCBFlag                bool
	Positions               []Position
	InitialBalances         InitialBalance
	CurrentBalances         CurrentBalance
	ProjectedBalances       ProjectedBalance
}

type Position struct {
	ShortQuantity                  float64
	AveragePrice                   float64
	CurrentDayProfitLoss           float64
	CurrentDayProfitLossPercentage float64
	LongQuantity                   float64
	SettledLongQuantity            float64
	SettledShortQuantity           float64
	AgedQuantity                   float64
	Instrument                     AccountInstrument
	MarketValue                    float64
	MaintenanceRequirement         float64
	AverageLongPrice               float64
	AverageShortPrice              float64
	TaxLotAverageLongPrice         float64
	TaxLotAverageShortPrice        float64
	LongOpenProfitLoss             float64
	ShortOpenProfitLoss            float64
	PreviousSessionLongQuantity    float64
	PreviousSessionShortQuantity   float64
	CurrentDayCost                 float64
}

type AccountInstrument struct {
	AssetType    string
	Cusip        string
	Symbol       string
	Description  string
	InstrumentID int
	NetChange    float64
}

type InitialBalance struct {
	AccruedInterest                  float64
	AvailableFundsNonMarginableTrade float64
	BondValue                        float64
	BuyingPower                      float64
	CashBalance                      float64
	CashAvailableForTrading          float64
	CashReceipts                     float64
	DayTradingBuyingPower            float64
	DayTradingBuyingPowerCall        float64
	DayTradingEquityCall             float64
	Equity                           float64
	EquityPercentage                 float64
	LiquidationValue                 float64
	LongMarginValue                  float64
	LongOptionMarketValue            float64
	LongStockValue                   float64
	MaintenanceCall                  float64
	MaintenanceRequirement           float64
	Margin                           float64
	MarginEquity                     float64
	MoneyMarketFund                  float64
	MutualFundValue                  float64
	RegTCall                         float64
	ShortMarginValue                 float64
	ShortOptionMarketValue           float64
	ShortStockValue                  float64
	TotalCash                        float64
	IsInCall                         bool
	UnsettledCash                    float64
	PendingDeposits                  float64
	MarginBalance                    float64
	ShortBalance                     float64
	AccountValue                     float64
}

/*type CurrentBalance struct {
	AvailableFunds                   float64
	AvailableFundsNonMarginableTrade float64
	BuyingPower                      float64
	BuyingPowerNonMarginableTrade    float64
	DayTradingBuyingPower            float64
	DayTradingBuyingPowerCall        float64
	Equity                           float64
	EquityPercentage                 float64
	LongMarginValue                  float64
	MaintenanceCall                  float64
	MaintenanceRequirement           float64
	MarginBalance                    float64
	RegTCall                         float64
	ShortBalance                     float64
	ShortMarginValue                 float64
	SMA                              float64
	IsInCall                         float64
	StockBuyingPower                 float64
	OptionBuyingPower                float64
}*/

type CurrentBalance struct {
	AccruedInterest                  float64
	CashBalance                      float64
	CashReceipts                     float64
	LongOptionMarketValue            float64
	LiquidationValue                 float64
	LongMarketValue                  float64
	MoneyMarketFund                  float64
	Savings                          float64
	ShortMarketValue                 float64
	PendingDeposits                  float64
	MutualFundValues                 float64
	BondValue                        float64
	ShortOptionMarketValue           float64
	AvailableFunds                   float64
	AvailableFundsNonMarginableTrade float64
	BuyingPower                      float64
	BuyingPowerNonMarginableTrade    float64
	DayTradingBuyingPower            float64
	Equity                           float64
	EquityPercentage                 float64
	LongMarginValue                  float64
	MaintenanceCall                  float64
	MaintenanceRequirement           float64
	MarginBalance                    float64
	RegTCall                         float64
	ShortBalance                     float64
	ShortMarginValue                 float64
	SMA                              float64
}

type ProjectedBalance struct {
	AvailableFunds                   float64
	AvailableFundsNonMarginableTrade float64
	BuyingPower                      float64
	BuyingPowerNonMarginableTrade    float64
	DayTradingBuyingPower            float64
	DayTradingBuyingPowerCall        float64
	Equity                           float64
	EquityPercentage                 float64
	LongMarginValue                  float64
	MaintenanceCall                  float64
	MaintenanceRequirement           float64
	MarginBalance                    float64
	RegTCall                         float64
	ShortBalance                     float64
	ShortMarginValue                 float64
	SMA                              float64
	IsInCall                         bool
	StockBuyingPower                 float64
	OptionBuyingPower                float64
}

type AggregatedBalance struct {
	CurrentLiquidationValue float64
	LiquidationValue        float64
}

type FullOrder struct {
	Session                  string
	Duration                 string
	OrderType                string
	CancelTime               string
	ComplexOrderStrategyType string
	Quantity                 int
	FilledQuantity           int
	RemainingQuantity        int
	RequestedDestination     string
	DestinationLinkName      string
	ReleaseTime              string
	StopPrice                int
	StopPriceLinkBasis       string
	StopPriceLinkType        string
	StopPriceOffset          int
	StopType                 string
	Price                    string
	TaxLotMethod             string
	OrderLegCollection       []FullOrderLeg
	ActivationPrice          int
	SpecialInstruction       string
	OrderStrategyType        string
	OrderId                  int
	Cancelable               bool
	Editable                 bool
	Status                   string
	EnteredTime              string
	CloseTime                string
	Tag                      string
	AccountNumber            int
	OrderActivityCollection  []FullOrderActivity
	ReplacingOrderCollection string
	ChildOrderStrategies     string
	StatusDescription        string
}

type FullOrderActivity struct {
	ActivityType           string
	ExecutionType          string
	Quantity               int
	OrderRemainingQuantity int
	ExecutionLegs          []FullExecutionLeg
}

type FullExecutionLeg struct {
	LegId             int
	Price             int
	Quantity          int
	MismarkedQuantity int
	InstrumentId      int
	Time              string
}

type FullOrderLeg struct {
	OrderLegType   string
	LegId          int
	Instrument     InstrumentRef
	Instruction    string
	PositionEffect string
	Quantity       int
	QuantityType   string
	DivCapGains    string
	ToSymbol       string
}

type SingleLegOrder struct {
	OrderType          string     `json:"orderType"`
	Session            string     `json:"session"`
	Duration           string     `json:"duration"`
	Strategy           string     `json:"orderStrategyType"`
	OrderLegCollection []OrderLeg `json:"orderLegCollection"`
}

type OrderLeg struct {
	Instruction string                `json:"instruction"`
	Quantity    int                   `json:"quantity"`
	Instrument  SimpleOrderInstrument `json:"instrument"`
}

type MultiLegOrder struct {
	OrderType          string // LIMIT, MARKET
	Session            string // NORMAL
	Duration           string // DAY
	Strategy           string // SINGLE
	OrderLegCollection []SimpleOrderLeg
}

type SimpleOrderLeg struct {
	Instruction string
	Quantity    int
	Instrument  SimpleOrderInstrument
}

type SimpleOrderInstrument struct {
	Symbol    string `json:"symbol"`
	AssetType string `json:"assetType"`
}

type (
	SingleLegOrderComposition      func(order *SingleLegOrder)
	MultiLegSimpleOrderComposition func(order *MultiLegOrder)
)
