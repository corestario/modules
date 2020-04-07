## Reseeding module

Implements reseeding logic for [RandApp](https://github.com/corestario/randapp) application.

### How it works

* Participants sends seeds to the module as a transaction using the `MsgSeed` message:
``` go
type MsgSeed struct {
	Sender sdk.AccAddress
	Seed   []byte
}
```
* Module stores those messages until you get a polka (2/3 + 1) votes for a given seed. Then the seed will be transferred to the Tendermint Consensus by calling the EndBlock() function
* Multiple seeds may be competing for a polka; as soon as one seed gets enough votes and gets transferred to Tendermint, all loser seeds are deleted.

### Example of seed sending

```shell script
./${YOUR_APP_CLI_BIN} tx reseeding send ${YOUR_SEED} --chain-id=${YOUR_CHAIN} --from ${ACCOUNT}
```

Note that there is a [reseeder](https://github.com/corestario/randapp/blob/master/reseeder.sh) script implemented by corestar.io that implements a bitcoin last-block-hash based seed. 
