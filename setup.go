package mapping

import (
  "context"
  "fmt"
  "net/url"
  "strconv"
  "strings"
  "time"

  "github.com/coredns/caddy"
  "github.com/coredns/coredns/core/dnsserver"
  "github.com/coredns/coredns/plugin"

  "github.com/go-redis/redis/v8"
)

func init() { plugin.Register("mapping", setup) }

func setup(c *caddy.Controller) error {
  // c.Next() // 'mapping'
  // if c.NextArg() {
  //   return plugin.Error("mapping", c.ArgErr())
  // }

  falcon, err := parseFalcon(c)

  if err != nil {
    return plugin.Error("mapping", err)
  }

  dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
    return Mapping{Next: next, FalconInstance: falcon}
  })

  return nil
}

const DEFAULT_CONNECT_TIMEOUT = time.Second * 10
const DEFAULT_READ_TIMEOUT = time.Second * 30
const DEFAULT_TTL = 300
const DEFAULT_PREFIX = "_dns_"
const DEFAULT_POSTFIX = ""

type Falcon struct {
  Url *url.URL
  Prefix string
  Suffix string
  Ttl uint32
  ConnectTimeout time.Duration
  ReadTimeout time.Duration
  RedisClient *redis.Client
}

func (falcon *Falcon) Connect() error {
  redisUrl := falcon.Url

  redisOptions := &redis.Options{
    Addr:     redisUrl.Host,
    DB: 0,
    DialTimeout: falcon.ConnectTimeout,
    ReadTimeout: falcon.ReadTimeout,
  }

  if redisUrl.User != nil {
    password, _ := redisUrl.User.Password()
    redisOptions.Password = password
  }
  

  if redisUrl.Path != "" {
    dbString := strings.Trim(redisUrl.Path, "/") 

    value, err := strconv.Atoi(dbString)

    if err != nil {
      return err
    }

    redisOptions.DB = value

  }

  redisClient := redis.NewClient(redisOptions)

  _, err := redisClient.Ping(context.Background()).Result()

  if err != nil {
    return err
  }

  fmt.Printf("Redis %s connected!\n", falcon.Url)

  falcon.RedisClient = redisClient

  return nil
}

func parseFalcon(c *caddy.Controller) (*Falcon, error) {
  falcon := Falcon {
    Url: &url.URL{
      Scheme: "redis",
      Host: "localhost:6379",
    },
    Prefix: DEFAULT_PREFIX,
    Suffix: DEFAULT_POSTFIX,
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

          if err != nil {
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

          if err != nil {
            return nil, c.ArgErr()
          }

          falcon.ConnectTimeout = time.Millisecond * time.Duration(value)

        case "read_timeout":
          if !c.NextArg() {
            return nil, c.ArgErr()
          }

          value, err := strconv.ParseInt(c.Val(), 10, 0)

          if err != nil {
            return nil, c.ArgErr()
          }

          falcon.ReadTimeout = time.Millisecond * time.Duration(value)

        case "ttl":
          if !c.NextArg() {
            return nil, c.ArgErr()
          }

          value, err := strconv.Atoi(c.Val())

          if err != nil {
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

  err = falcon.Connect()

  if err != nil {
    return nil, err
  }

  return &falcon, nil
}
