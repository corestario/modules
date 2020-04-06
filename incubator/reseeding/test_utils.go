package reseeding

import (
	"bytes"
	"encoding/hex"
	"strconv"

	"github.com/corestario/modules/incubator/reseeding/internal/keeper"
	"github.com/corestario/modules/incubator/reseeding/internal/types"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

var maccPerms = map[string][]string{
	auth.FeeCollectorName:     nil,
	distr.ModuleName:          nil,
	staking.BondedPoolName:    {supply.Burner, supply.Staking},
	staking.NotBondedPoolName: {supply.Burner, supply.Staking},
	types.ModuleName:          nil,
}

var ModuleBasics = module.NewBasicManager(
	genutil.AppModuleBasic{},
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	params.AppModuleBasic{},
	staking.AppModuleBasic{},
	distr.AppModuleBasic{},
	slashing.AppModuleBasic{},
	supply.AppModuleBasic{},
	AppModule{},
)

type testApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	keyMain      *sdk.KVStoreKey
	keyAccount   *sdk.KVStoreKey
	keySupply    *sdk.KVStoreKey
	keyStaking   *sdk.KVStoreKey
	tkeyStaking  *sdk.TransientStoreKey
	keyDistr     *sdk.KVStoreKey
	tkeyDistr    *sdk.TransientStoreKey
	keyNFT       *sdk.KVStoreKey
	keyParams    *sdk.KVStoreKey
	tkeyParams   *sdk.TransientStoreKey
	keySlashing  *sdk.KVStoreKey
	keyReseeding *sdk.KVStoreKey

	// Keepers
	accountKeeper   auth.AccountKeeper
	bankKeeper      bank.Keeper
	supplyKeeper    supply.Keeper
	stakingKeeper   staking.Keeper
	slashingKeeper  slashing.Keeper
	distrKeeper     distr.Keeper
	paramsKeeper    params.Keeper
	reseedingKeeper keeper.Keeper

	// Module Manager
	mm *module.Manager
}

func getTestEnvironment() (sdk.Context, keeper.Keeper, staking.Keeper, []crypto.PubKey) {
	cdc := types.ModuleCdc

	bApp := bam.NewBaseApp("test", log.NewNopLogger(), dbm.NewMemDB(), auth.DefaultTxDecoder(cdc))

	// Here you initialize your application with the store keys it requires
	var app = &testApp{
		BaseApp: bApp,
		cdc:     cdc,

		keyMain:      sdk.NewKVStoreKey(bam.MainStoreKey),
		keyAccount:   sdk.NewKVStoreKey(auth.StoreKey),
		keySupply:    sdk.NewKVStoreKey(supply.StoreKey),
		keyStaking:   sdk.NewKVStoreKey(staking.StoreKey),
		tkeyStaking:  sdk.NewTransientStoreKey(staking.TStoreKey),
		keyDistr:     sdk.NewKVStoreKey(distr.StoreKey),
		tkeyDistr:    sdk.NewTransientStoreKey("transient_" + distr.ModuleName),
		keyParams:    sdk.NewKVStoreKey(params.StoreKey),
		tkeyParams:   sdk.NewTransientStoreKey(params.TStoreKey),
		keySlashing:  sdk.NewKVStoreKey(slashing.StoreKey),
		keyReseeding: sdk.NewKVStoreKey(types.StoreKey),
	}

	app.paramsKeeper = params.NewKeeper(app.cdc, app.keyParams, app.tkeyParams)
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount,
		authSubspace,
		auth.ProtoBaseAccount,
	)

	// The BankKeeper allows you perform sdk.Coins interactions
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		bankSubspace,
		app.ModuleAccountAddrs(),
	)

	app.supplyKeeper = supply.NewKeeper(app.cdc, app.keySupply, app.accountKeeper,
		app.bankKeeper, maccPerms)

	// The staking keeper
	// The staking keeper
	app.stakingKeeper = staking.NewKeeper(
		app.cdc,
		app.keyStaking,
		app.supplyKeeper,
		stakingSubspace,
	)

	reseedingKeeper := keeper.NewKeeper(app.cdc, app.keyReseeding)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
	)

	app.mm.SetOrderBeginBlockers(distr.ModuleName, slashing.ModuleName)
	app.mm.SetOrderEndBlockers(staking.ModuleName)

	// Sets the orandappder of Genesis - Orandappder matters, genutil is to always come last
	app.mm.SetOrderInitGenesis(
		distr.ModuleName,
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,

		types.ModuleName,

		genutil.ModuleName,
	)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(
		auth.NewAnteHandler(
			app.accountKeeper, app.supplyKeeper, auth.DefaultSigVerificationGasConsumer,
		),
	)

	app.MountStores(
		app.keyMain,
		app.keyAccount,
		app.keySupply,
		app.keyStaking,
		app.tkeyStaking,
		app.keyDistr,
		app.tkeyDistr,
		app.keySlashing,
		app.keyParams,
		app.tkeyParams,
		app.keyReseeding,
	)

	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		panic(err)
	}

	genesisState := ModuleBasics.DefaultGenesis()
	stateBytes, err := codec.MarshalJSONIndent(app.cdc, genesisState)
	if err != nil {
		panic(err)
	}

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			AppStateBytes: stateBytes,
		},
	)

	ctx := app.NewContext(false, abci.Header{})

	PKs := createTestPubKeys(4)

	for _, pk := range PKs {
		valPubKey := pk
		valAddr := sdk.ValAddress(valPubKey.Address().Bytes())
		valTokens := sdk.TokensFromConsensusPower(10)
		validator := stakingTypes.NewValidator(valAddr, valPubKey, stakingTypes.Description{})
		validator, _ = validator.AddTokensFromDel(valTokens)
		app.stakingKeeper.SetValidator(ctx, validator)
	}

	return ctx, reseedingKeeper, app.stakingKeeper, PKs
}

func createTestPubKeys(numPubKeys int) []crypto.PubKey {
	var publicKeys []crypto.PubKey
	var buffer bytes.Buffer

	//start at 10 to avoid changing 1 to 01, 2 to 02, etc
	for i := 100; i < (numPubKeys + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF") //base pubkey string
		buffer.WriteString(numString)                                                       //adding on final two digits to make pubkeys unique
		publicKeys = append(publicKeys, NewPubKey(buffer.String()))
		buffer.Reset()
	}
	return publicKeys
}

func NewPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	//res, err = crypto.PubKeyFromBytes(pkBytes)
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes)
	return pkEd
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *testApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}
