package conf

import (
	"flag"
	"os"

	"github.com/go-kratos/kratos/contrib/config/apollo/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	_ "github.com/go-kratos/kratos/v2/encoding/json"
	_ "github.com/go-kratos/kratos/v2/encoding/yaml"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	c config.Config

	// flagconf is the config flag.
	flagconf string

	// apollo
	APOLLO_APPID   = "1021"
	APOLLO_CLUSTER = os.Getenv("APP_ENV")
	APOLLO_URL     = os.Getenv("APOLLO_URL")
	APOLLO_NS      = "bootstrap.yaml"
	APOLLO_SECRET  = os.Getenv("APOLLO_SECRET")
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	// TODO:
	flag.Parse()

	// file conf
	var cfs = file.NewSource(flagconf)
	// apollo conf
	if _, ok := os.LookupEnv("APOLLO_URL"); ok {
		cfs = apollo.NewSource(
			apollo.WithAppID(APOLLO_APPID),
			apollo.WithCluster(APOLLO_CLUSTER),
			apollo.WithEndpoint(APOLLO_URL),
			apollo.WithNamespace(APOLLO_NS),
			apollo.WithSecret(APOLLO_SECRET),
			apollo.WithEnableBackup(),
		)
		// debug
		log.Infof("[config] init with apollo: APPID:%s, CLUSTER:%s, URL:%s, NS:%s, SECRET:%s\n", APOLLO_APPID, APOLLO_CLUSTER, APOLLO_URL, APOLLO_NS, APOLLO_SECRET)
	}

	c = config.New(config.WithSource(cfs))
	if err := c.Load(); err != nil {
		panic(err)
	}
}

// get bootstrap conf
func GetConf() (*Bootstrap, error) {
	var bc Bootstrap

	if err := c.Scan(&bc); err != nil {
		return nil, err
	}

	return &bc, nil
}

// example: GetKey("bootstrap.data", &data)
func GetKey(key string, val interface{}) error {
	return c.Value(key).Scan(val)
}
