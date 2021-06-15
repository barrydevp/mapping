package mapping

import (
  "time"
  "net/url"
  "strconv"

  "github.com/coredns/caddy"
  "github.com/coredns/coredns/core/dnsserver"
  "github.com/coredns/coredns/plugin"

  "github.com/go-redis/redis/v8"
)

func init() { plugin.Register("mapping", setup) }

func setup(c *caddy.Controller) error {
  c.Next() // 'mapping'
  if c.NextArg() {
    return plugin.Error("mapping", c.ArgErr())
  }

  dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
    return Mapping{Next: next}
  })

  return nil
}

type Falcon struct {
  Url *url.URL
  Prefix string
  Suffix string
  Ttl uint32
  ConnectTimeout time.Duration
  ReadTimeout time.Duration
  RedisClient redis.Client
}

func (falcon *Falcon) Connect() {
  redisUrl := falcon.Url
  
  falcon.RedisClient = redis.NewClient(&redis.Options{
    Addr:     redisUrl.Host,
    Password: redisUrl.User.Password, // no password set
    DB:       0,  // use default DB
  })
}


const DEFAULT_CONNECT_TIMEOUT = time.Duration(30000)
const DEFAULT_READ_TIMEOUT = time.Duration(30000)
const DEFAULT_TTL = 300

func parseFalcon(c *caddy.Controller) (*Falcon, error) {
  falcon := Falcon {
    Url: &url.URL{
      Scheme: "redis",
      Host: "localhost:6379",
    },
    Prefix: "_dns",
    Suffix: "",
    Ttl: DEFAULT_TTL,
    ConnectTimeout: DEFAULT_CONNECT_TIMEOUT,
    ReadTimeout: DEFAULT_READ_TIMEOUT,
  }

  var err error

  for c.Next() {
    if c.NextBlock() {
      for {
        switch c.Val() {
        case "uri":
          if !c.NextArg() {
            return nil, c.ArgErr()
          }

          falcon.Url, err = url.Parse(c.Val())

          if(err != nil) {
            return nil, c.ArgErr()
          }

        case "prefix":
          if !c.NextArg() {
            return nil, c.ArgErr()
          }
          
          falcon.Prefix = c.Val()

        case "suffix":
          if !c.NextArg() {
            return nil, c.ArgErr()
          }

          falcon.Suffix = c.Val()

        case "connect_timeout":
          if !c.NextArg() {
            return nil, c.ArgErr()
          }
          
          value, err := strconv.ParseInt(c.Val(), 10, 0)

          if(err != nil) {
            return nil, c.ArgErr()
          }

          falcon.ConnectTimeout = time.Duration(value)

        case "read_timeout":
          if !c.NextArg() {
            return nil, c.ArgErr()
          }

          value, err := strconv.ParseInt(c.Val(), 10, 0)

          if(err != nil) {
            return nil, c.ArgErr()
          }

          falcon.ReadTimeout = time.Duration(value)

        case "ttl":
          if !c.NextArg() {
            return nil, c.ArgErr()
          }

          value, err := strconv.Atoi(c.Val())

          if(err != nil) {
            return nil, c.ArgErr()
          }

          falcon.Ttl = uint32(value)

        default:
          if c.Val() != "}" {
            return nil, c.Errf("unknown property '%s'", c.Val())
          }
        }

        if !c.Next() {
          break
        }
      }
    }
  }

  return &falcon, nil
}
