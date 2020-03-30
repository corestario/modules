module github.com/corestario/modules/incubator/reseeding

go 1.12

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20191013030331-92ea174ea6e6
	github.com/gorilla/mux v1.7.3
	github.com/spf13/cobra v0.0.5
	github.com/tendermint/tendermint v0.32.8
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/corestario/cosmos-sdk v0.3.0
	github.com/tendermint/tendermint => ./../../../tendermint
	go.dedis.ch/kyber/v3 => github.com/corestario/kyber/v3 v3.0.0-20200218082721-8ed10c357c05
)
