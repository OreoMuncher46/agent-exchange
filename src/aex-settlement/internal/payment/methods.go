package payment

// SettlementMethod describes one AP2 Network & Issuer option available to
// the Settlement service. The mandate/credential flow stays intact; only the
// network underneath changes.
type SettlementMethod struct {
	ID string `json:"id"` // e.g. "card:usd", "nano:mainnet"
	// Network is the settlement network. "nano:mainnet" is native Nano.
	Network string `json:"network"`
	// Asset settled on the network, e.g. "USD", "USDC", "XNO".
	Asset string `json:"asset"`
	// DisplayName is human-readable, e.g. "Nano (XNO)".
	DisplayName string `json:"display_name"`
	// FeeFixedUSD is the per-transaction fee floor in USD (gas, bridge,
	// processor minimum). Nano has no fee floor.
	FeeFixedUSD float64 `json:"fee_fixed_usd"`
	// FeePercent is the proportional fee (e.g. 0.029 for 2.9%).
	FeePercent float64 `json:"fee_percent"`
	// MinViableBountyUSD is the smallest bounty that stays viable on this
	// network after fees. Nano stays viable at any size.
	MinViableBountyUSD float64 `json:"min_viable_bounty_usd"`
	// Finality describes settlement speed, e.g. "sub-second".
	Finality string `json:"finality"`
	// RequiresIssuer is true when settlement depends on a card issuer.
	RequiresIssuer bool `json:"requires_issuer"`
	// RequiresBridge is true when settlement needs a wrapped asset or bridge.
	RequiresBridge bool `json:"requires_bridge"`
}

const (
	// MethodCardUSD is card/USD settlement via MPP -> Network & Issuer.
	MethodCardUSD = "card:usd"
	// MethodBaseUSDC is Base-USDC settlement (gas + bridge floor applies).
	MethodBaseUSDC = "base:usdc"
	// MethodNanoMainnet is feeless Nano (XNO) settlement for micro-bounties.
	MethodNanoMainnet = "nano:mainnet"
	// MicroBountyThresholdUSD is the fee-floor boundary: bounties below this
	// SHOULD settle via nano:mainnet so a $0.50 payout is not consumed by
	// gas + bridge/per-transaction minimums.
	MicroBountyThresholdUSD = 1.0
)

// settlementMethodTable is the Settlement-service payment-method table.
var settlementMethodTable = []SettlementMethod{
	{
		ID:                 MethodCardUSD,
		Network:            "card:usd",
		Asset:              "USD",
		DisplayName:        "Card (USD)",
		FeeFixedUSD:        0.30,
		FeePercent:         0.029,
		MinViableBountyUSD: 1.0,
		Finality:           "days",
		RequiresIssuer:     true,
		RequiresBridge:     false,
	},
	{
		ID:                 MethodBaseUSDC,
		Network:            "base:usdc",
		Asset:              "USDC",
		DisplayName:        "Base USDC",
		FeeFixedUSD:        0.05,
		FeePercent:         0.0,
		MinViableBountyUSD: 0.5,
		Finality:           "seconds",
		RequiresIssuer:     false,
		RequiresBridge:     true,
	},
	{
		ID:                 MethodNanoMainnet,
		Network:            "nano:mainnet",
		Asset:              "XNO",
		DisplayName:        "Nano (XNO)",
		FeeFixedUSD:        0.0,
		FeePercent:         0.0,
		MinViableBountyUSD: 0.0,
		Finality:           "sub-second",
		RequiresIssuer:     false,
		RequiresBridge:     false,
	},
}

// ListSettlementMethods returns the Settlement-service payment-method table.
func ListSettlementMethods() []SettlementMethod {
	methods := make([]SettlementMethod, len(settlementMethodTable))
	copy(methods, settlementMethodTable)
	return methods
}

// GetSettlementMethod returns the method with the given ID.
func GetSettlementMethod(id string) (SettlementMethod, bool) {
	for _, m := range settlementMethodTable {
		if m.ID == id {
			return m, true
		}
	}
	return SettlementMethod{}, false
}

// EstimateSettlementFeeUSD estimates the network fee for amountUSD in USD.
// Nano (XNO) is feeless: it always returns 0.
func EstimateSettlementFeeUSD(methodID string, amountUSD float64) float64 {
	method, ok := GetSettlementMethod(methodID)
	if !ok {
		return 0
	}
	if amountUSD < 0 {
		amountUSD = 0
	}
	return method.FeeFixedUSD + amountUSD*method.FeePercent
}

// SelectSettlementMethodForAmount picks the settlement network for a bounty.
// Amounts below the fee floor (MicroBountyThresholdUSD) select nano:mainnet
// so sub-dollar payouts stay viable; larger amounts use card:usd.
func SelectSettlementMethodForAmount(amountUSD float64) SettlementMethod {
	if amountUSD < MicroBountyThresholdUSD {
		if m, ok := GetSettlementMethod(MethodNanoMainnet); ok {
			return m
		}
	}
	m, _ := GetSettlementMethod(MethodCardUSD)
	return m
}
