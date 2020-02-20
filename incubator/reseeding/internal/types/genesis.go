package types

// GenesisState is the state that must be provided at genesis.
type GenesisState struct {
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{}
}

// DefaultGenesisState returns a default genesis state.
func DefaultGenesisState() GenesisState {
	return NewGenesisState()
}

func ValidateGenesis(data GenesisState) error {
	return nil
}
