# State

## Seeds

`Seeds` is a collection of seeds that are mapped to the votes that those seeds received from validators.

```go
// Seeds is a map from a seed to its senders.
type Seeds map[string]SeedVotes
```

You can add votes to currently stored seeds:
```go
func (s Seeds) Add(seed []byte, sender string) {
	key := GetSeedKey(seed)
	votes, ok := s[key]
	if !ok {
		votes = make(map[string]struct{})
	}
	votes[sender] = struct{}{}
	s[key] = votes
}
```