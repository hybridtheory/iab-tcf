# iab-tcf [![CircleCI](https://circleci.com/gh/affectv/iab-tcf.svg?style=svg)](https://circleci.com/gh/affectv/iab-tcf)

Go code to parse IAB v1/v2 consent based on [github.com/LiveRamp/iabconsent](github.com/LiveRamp/iabconsent) library. 

This package is just a wrapper on their amazing work to fill our needs.

## Usage

### Install

```bash
go get github.com/affectv/iab-tcf
```

### Example

```golang
package main

import (
    "fmt"
    iab "github.com/affectv/iab-tcf"
)

func main() {
    encoded := "COyt4MbOyt4MbMOAAAENAiCgAIAAAAAAAAAAADEAAgIAAAAAAAA.YAAAAAAAAAA"
    consent, err := iab.NewConsent(encoded)
}
```

## Testing

We use [Ginkgo](https://onsi.github.io/ginkgo/) and [Gomega](https://onsi.github.io/gomega/) for testing purposes.

```bash
ginkgo ./...
```
