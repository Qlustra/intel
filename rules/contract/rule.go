package contract

type StructContractRule struct {
	RuleType
}

func Contract() *StructContractRule {
	return &StructContractRule{}
}
