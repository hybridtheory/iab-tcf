package cmp

type Consent struct{}

// ValidCMPs returns the list of valid CMPs loaded.
func (c *Consent) ValidCMPs() []int {
	return ValidCMPs
}

// IsCMPListLoaded returns if the list of valid CMPs was loaded or not.
func (c *Consent) IsCMPListLoaded() bool {
	return c.ValidCMPs() != nil
}
