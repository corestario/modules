module github.com/cosmos/modules/incubator/reseeding

go 1.13

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20191013030331-92ea174ea6e6
	github.com/gogo/protobuf v1.3.0
	github.com/gorilla/mux v1.7.3
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.1
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/tendermint v0.32.8
	github.com/tendermint/tm-db v0.3.0
)

replace (
	github.com/corestario/cosmos-utils/client => ./../../../cosmos-utils/client
	github.com/corestario/dkglib => ./../../../dkglib
	github.com/cosmos/cosmos-sdk => ./../../../cosmos-sdk
	github.com/tendermint/tendermint => ./../../../tendermint
	go.dedis.ch/kyber/v3 => ./../../../kyber
)
