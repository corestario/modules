# Messages

## MsgSeed 

`MsgSeed` message is the only message used by the module. It holds the seed and sender address. 

| **Field** | **Type**         | **Description**                                              |
|:----------|:-----------------|:-------------------------------------------------------------|
| Sender    | `sdk.AccAddress` | The account address of the user sending the seed.            |
| Seed      | `[]byte`         | The seed (an arbitrary sequence of bytes. Max size is 1 MiB. |

``` go
type MsgSeed struct {
    Sender sdk.AccAddress
    Seed   []byte
}
```
