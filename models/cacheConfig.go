package models

type CacheConfig struct {
	CacheURL            string
	CacheEnabled        bool
	CacheTTL            string
	CacheCaCertPath     string
	CacheClientCertPath string
	CacheClientKeyPath  string
}
