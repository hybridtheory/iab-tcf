# iab-tcf [![CircleCI](https://circleci.com/gh/hybridtheory/iab-tcf.svg?style=svg)](https://circleci.com/gh/hybridtheory/iab-tcf)

Go code to parse IAB v1/v2 consent based on [github.com/LiveRamp/iabconsent](github.com/LiveRamp/iabconsent)
library.

This package is just a wrapper on their amazing work to fill our needs.

## Usage

### Install

```bash
go get github.com/hybridtheory/iab-tcf
```

### Example

```golang
package main

import (
    "fmt"
    iab "github.com/hybridtheory/iab-tcf"
)

func main() {
    encoded := "COyt4MbOyt4MbMOAAAENAiCgAIAAAAAAAAAAADEAAgIAAAAAAAA.YAAAAAAAAAA"
    consent, err := iab.NewConsent(encoded)
}
```

### CMP vendor list

In order to validate if the CMP received belongs to a valid id there's a loader integrated
with the library.

```golang
package main

import (
    "fmt"
    iab "github.com/hybridtheory/iab-tcf"
    cmp "github.com/hybridtheory/iab-tcf/cmp"
)

func main() {
    cmp.NewLoader().LoadIDs()
    consent, err := iab.NewConsent(encoded)
    consent.IsCMPValid() // This will return true/false depending on the ids loaded.
}
```

If instead of the default vendor list [https://cmplist.consensu.org/v2/cmp-list.json](https://cmplist.consensu.org/v2/cmp-list.json)
we want to use our own, we can use an option:

```golang
err := cmp.NewLoader(cmp.WithURL("https://example.com/cmp-list.json")).LoadIDs()
```

The format of the JSON must be the same.

If we want to use our own list of valid CMPs we can simply set the variable:

```golang
cmp.ValidCMPs = []int{1, 123}
```

## Testing

We use [Ginkgo](https://onsi.github.io/ginkgo/) and [Gomega](https://onsi.github.io/gomega/)
for testing purposes.

```bash
ginkgo ./...
```
