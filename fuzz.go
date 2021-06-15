// +build gofuzz

package mapping

import (
  "github.com/coredns/coredns/plugin/pkg/fuzz"
)

// Fuzz fuzzes cache.
func Fuzz(data []byte) int {
  m := Mapping{}
  return fuzz.Do(m, data)
}
