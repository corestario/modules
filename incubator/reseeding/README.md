## Reseeding module

Implements reseeding logic for [RandApp](https://github.com/corestario/randapp) application.

### How it works

* Participants sends seeds to the module
* Module stores them
* When count one of the seeds is greater or equal to 2/3*ValidatorsCount, the seed will be transferred to the Tendermint Consensus by calling the EndBlock() function

### Example of seed sending

```shell script
YOUR_APP_CLI_BIN reseeding send "SEED_BYTES"
```