package types

import "fmt"

func (p *IBCParams) String() string {
	return fmt.Sprintf(`Params:
  Unbonding Time:     %s
  Max Validators:     %d
  Max Entries:        %d
  Historical Entries: %d
  Bonded Coin Denom:  %s`, p.UnbondingTime,
		p.MaxValidators, p.MaxEntries, p.HistoricalEntries, p.BondDenom)
}
