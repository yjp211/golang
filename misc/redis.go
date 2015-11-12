package misc

import (
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/garyburd/redigo/redis"
)

type PoolProxy struct {
	*redis.Pool
	Addr string
}

type RedisCluster struct {
	Proxys  []*PoolProxy
	Count   int
	Counter int64
	Debug   bool
	Logger  *log.Logger
}

type Conn struct {
	redis.Conn
	proxy   *PoolProxy
	cluster *RedisCluster
}

func (p *RedisCluster) GetConn() Conn {
	var proxy *PoolProxy
	if p.Count == 1 {
		proxy = p.Proxys[0]
	} else {
		counter := p.inrcCounter()
		if counter >= 100000 {
			p.resetCounter()
			counter = 0
		}
		index := int(counter) % p.Count
		proxy = p.Proxys[index]

	}
	raw_conn := proxy.Get()

	return Conn{
		Conn:    raw_conn,
		cluster: p,
		proxy:   proxy,
	}
}

func (p *RedisCluster) inrcCounter() int64 {
	return atomic.AddInt64(&p.Counter, 1)
}
func (p *RedisCluster) resetCounter() {
	atomic.StoreInt64(&p.Counter, 0)
}

func (p *RedisCluster) SetLogger(logger *log.Logger) {
	p.Logger = logger
}

func (p Conn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if p.cluster.Debug {
		start := time.Now()
		reply, err = p.Conn.Do(commandName, args...)
		end := time.Now()
		use := end.Sub(start)
		if err == nil {
			p.cluster.Logger.Printf("Proxy<%s> command<%s %v>, use: %v, reply:%v\n", p.proxy.Addr, commandName, args, use, reply)
		} else {
			p.cluster.Logger.Printf("Proxy<%s> command<%s %v>, use: %v, error:%v\n,", p.proxy.Addr, commandName, args, use, err)
		}
	} else {
		reply, err = p.Conn.Do(commandName, args...)
	}
	return reply, err
}

func newPoolProxy(addr string, password string, max int, timeout int, wait bool) *PoolProxy {
	raw_pool := &redis.Pool{
		Wait:        wait,
		MaxIdle:     max,
		IdleTimeout: time.Second * time.Duration(timeout),
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return &PoolProxy{
		Pool: raw_pool,
		Addr: addr,
	}
}

func NewRedisCluster(addrs []string, password string, max int, timeout int) (*RedisCluster, error) {
	RedisCluster := &RedisCluster{
		Debug:  true,
		Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
	count := len(addrs)
	RedisCluster.Count = count
	RedisCluster.Counter = 0
	RedisCluster.Proxys = make([]*PoolProxy, count)
	for i := 0; i < count; i++ {
		RedisCluster.Proxys[i] = newPoolProxy(addrs[i], password, max, timeout, false)
	}
	return RedisCluster, nil
}
