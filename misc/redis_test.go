package misc

import (
	"testing"
)

func TestRedis(t *testing.T) {
	//path := "/tmp/log"
	//level := "DEBUG"
	//Log, err := NewLog(path, level)
	//if err != nil {
	//t.Errorf("create log  at<%s> error: %s\n", path, err)
	//}
	addrs := []string{"127.0.0.1:6379", "localhost:6479"}
	passwd := ""
	maxconn := 200
	timeout := 3

	Cluster, err := NewRedisCluster(addrs, passwd, maxconn, timeout)
	if err != nil {
		t.Errorf("create redis cluster  %s error, %s\n", addrs, err)
	}
	conn1 := Cluster.GetConn()
	defer conn1.Close()
	conn2 := Cluster.GetConn()
	defer conn2.Close()
	conn3 := Cluster.GetConn()
	defer conn3.Close()
	conn1.Do("SET", "hellp", "world")
	conn2.Do("HMSET", "xxxx", "a", 1, "b", "bbb")
	conn3.Do("HMGET", "xxxx", "a", "c", "b", "bbb")
}
