package types

type Config struct {
	RedisUrl       string
	RedisAuth      string
	RedisDB        int
	MysqlUrl       string
	MysqlUser      string
	MysqlPass      string
	CacheQueueSize int
}
