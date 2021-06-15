package mapping

import (
  "context"
  "errors"
  "fmt"

  "github.com/coredns/coredns/plugin"
  "github.com/coredns/coredns/plugin/pkg/log"

  "github.com/miekg/dns"
)

const name = "mapping"

type Mapping struct {
  Next plugin.Handler
  FalconInstance *Falcon
}

// Name implements the Handler interface.
func (m Mapping) Name() string { return name }

// ServeDNS implements the plugin.Handler interface.
func (m Mapping) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
  wr, err := m.mapDomain(w, r)

  if err != nil {
    log.Infof("[mapping]: cannot map domain, %s\n", err)
  }

  return plugin.NextOrFailure(m.Name(), m.Next, ctx, wr, r)
}

func (m *Mapping) mapDomain(w dns.ResponseWriter, r *dns.Msg) (dns.ResponseWriter, error) {

  originalQuestion := r.Question[0]

  falcon := m.FalconInstance
  redisClient := falcon.RedisClient

  domainQuery := fmt.Sprintf("%s%s%s", falcon.Prefix, originalQuestion.Name, falcon.Suffix)

  serviceDNS, err := redisClient.Get(context.Background(), domainQuery).Result()

  if err != nil {
    return w, err
  }
  
  if serviceDNS == "" {
    return w, errors.New("domain not found")
  }

  log.Infof("[mapping]: found %s\n", serviceDNS)

  if serviceDNS[len(serviceDNS) - 1] != '.' {
    serviceDNS += "."
  }

  r.Question[0].Name = serviceDNS

  return &MapperWriter{ResponseWriter: w, originalQuestion: originalQuestion, skip: false}, nil
}

