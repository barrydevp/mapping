package mapping

import (
  "github.com/miekg/dns"
)

type MapperWriter struct {
  dns.ResponseWriter
  originalQuestion dns.Question
  skip bool
}

func (m *MapperWriter) WriteMsg(originalRes *dns.Msg) error {
  if m.skip {
    return m.ResponseWriter.WriteMsg(originalRes)
  }

  // Deep copy 'res' as to not (e.g). rewrite a message that's also stored in the cache.
  res := originalRes.Copy()

  res.Question[0] = m.originalQuestion

  return m.ResponseWriter.WriteMsg(res)
}

