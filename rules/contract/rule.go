package contract

// StructContractRule is a composable rule for an embedded contract encounter.
type StructContractRule struct {
	RuleType
}

// Contract creates an unconstrained contract rule for adding expectations.
func Contract() *StructContractRule {
	return &StructContractRule{}
}
