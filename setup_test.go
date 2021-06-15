package mapping

import (
  "testing"

  "github.com/coredns/caddy"
)

func TestSetup(t *testing.T) {
  c := caddy.NewTestController("dns", `mapping`)
  if err := setup(c); err != nil {
    t.Fatalf("Expected no errors, but got: %v", err)
  }

  c = caddy.NewTestController("dns", `mapping example.org`)
  if err := setup(c); err == nil {
    t.Fatalf("Expected errors, but got: %v", err)
  }
}
