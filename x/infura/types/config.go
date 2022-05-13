package types

type Config struct {
	RedisUrl       string
	RedisAuth      string
	RedisDB        int
	MysqlUrl       string
	MysqlUser      string
	MysqlPass      string
	MysqlDB        string
	CacheQueueSize int
}
