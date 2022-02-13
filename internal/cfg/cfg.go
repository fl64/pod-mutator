package cfg

import (
	"github.com/ilyakaznacheev/cleanenv"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Resources struct {
	CPU string `yaml:"CPU"`
	MEM string `yaml:"MEM"`
}

type ReqLim struct {
	Limits   Resources `yaml:"limits"`
	Requests Resources `yaml:"requests"`
}

type MutatorConfig struct {
	Default  ReqLim           `yaml:"default"`
	Override []ReqLimOverride `yaml:"override"`
}

type ReqLimOverride struct {
	ImagePattern string `yaml:"image-pattern"`
	Resources    ReqLim `yaml:"resources"`
}

type LoggerCfg struct {
	DevMode bool `yaml:"dev" env:"DEV_MODE" env-default:"false"`
}

type Cfg struct {
	MetricAddr    string        `yaml:"metrics-bind-address" env:"METRICS_ADDR" env-default:":8080"`
	ProbeAddr     string        `yaml:"health-probe-bind-address" env:"PROBE_ADDR" env-default:":8081"`
	WebhookAddr   string        `yaml:"webhook-bind-address" env:"WEBHOOK_ADDR" env-default:":9443"`
	LeaderElect   bool          `yaml:"leader-elect" env:"LEADER_ELECT" env-default:"false"`
	LoggerCfg     LoggerCfg     `yaml:"log" env-prefix:"LOG_"`
	LabelSelector []string      `yaml:"label-selector" env:"LABEL_SELECTOR"`
	MutatorConfig MutatorConfig `yaml:"mutator-config"`
}

func GetCfg(configPath string) (*Cfg, error) {
	var cfg Cfg
	logger := ctrl.Log.WithName("config")
	logger.Info("Getting config form file")
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		logger.Info("Can't load config... skipping")
	}
	err = cleanenv.ReadEnv(&cfg)
	return &cfg, err
}
