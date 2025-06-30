// SPDX-License-Identifier: PolyForm-Noncommercial-1.0.0
// License: PolyForm Noncommercial License 1.0.0 (https://polyformproject.org/licenses/noncommercial/1.0.0/)

package main

import (
	"ai-rules-link/cmd"
	"embed"
)

//go:embed all:rules
var embeddedRules embed.FS

func main() {
	cmd.SetEmbeddedRules(embeddedRules)
	cmd.Execute(embeddedRules)
}
