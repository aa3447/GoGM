package gameLogic

import (
	"fmt"
)

type Buff struct{
	Attribute string
	Amount int
	Duration int // Duration in seconds
}

func (b *Buff) String() string {
	return fmt.Sprintf("%s: +%d for %d seconds", b.Attribute, b.Amount, b.Duration)
}
