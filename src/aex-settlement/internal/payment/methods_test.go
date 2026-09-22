package payment

import (
	"testing"
)

func TestSettlementMethodTableHasNanoMainnet(t *testing.T) {
	m, ok := GetSettlementMethod(MethodNanoMainnet)
	if !ok {
		t.Fatalf("GetSettlementMethod(%q) not found", MethodNanoMainnet)
	}
	if m.Network != "nano:mainnet" {
		t.Errorf("network = %q, want %q", m.Network, "nano:mainnet")
	}
	if m.Asset != "XNO" {
		t.Errorf("asset = %q, want %q", m.Asset, "XNO")
	}
	if m.FeeFixedUSD != 0 || m.FeePercent != 0 {
		t.Errorf("nano fees = (%v, %v), want feeless (0, 0)", m.FeeFixedUSD, m.FeePercent)
	}
	if m.RequiresIssuer || m.RequiresBridge {
		t.Errorf("nano requires_issuer=%v requires_bridge=%v, want (false, false)",
			m.RequiresIssuer, m.RequiresBridge)
	}
}

func TestEstimateSettlementFeeFloorComparison(t *testing.T) {
	// A $0.50 micro-bounty: card and Base-USDC floors dominate, Nano is free.
	const bounty = 0.50

	nanoFee := EstimateSettlementFeeUSD(MethodNanoMainnet, bounty)
	if nanoFee != 0 {
		t.Errorf("nano fee for $%.2f = %v, want 0", bounty, nanoFee)
	}

	cardFee := EstimateSettlementFeeUSD(MethodCardUSD, bounty)
	if cardFee <= nanoFee {
		t.Errorf("card fee %v should exceed nano fee %v for $%.2f bounty", cardFee, nanoFee, bounty)
	}

	baseFee := EstimateSettlementFeeUSD(MethodBaseUSDC, bounty)
	if baseFee <= nanoFee {
		t.Errorf("base-usdc fee %v should exceed nano fee %v for $%.2f bounty", baseFee, nanoFee, bounty)
	}
}

func TestSelectSettlementMethodForAmount(t *testing.T) {
	micro := SelectSettlementMethodForAmount(0.50)
	if micro.ID != MethodNanoMainnet {
		t.Errorf("select(0.50) = %q, want %q", micro.ID, MethodNanoMainnet)
	}

	standard := SelectSettlementMethodForAmount(150.00)
	if standard.ID != MethodCardUSD {
		t.Errorf("select(150.00) = %q, want %q", standard.ID, MethodCardUSD)
	}
}

func TestListSettlementMethodsNotEmpty(t *testing.T) {
	methods := ListSettlementMethods()
	if len(methods) == 0 {
		t.Fatal("ListSettlementMethods() returned no methods")
	}
	found := false
	for _, m := range methods {
		if m.ID == MethodNanoMainnet {
			found = true
		}
	}
	if !found {
		t.Errorf("ListSettlementMethods() missing %q entry", MethodNanoMainnet)
	}
}
