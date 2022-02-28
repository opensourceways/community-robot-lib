package options

import (
	"flag"
	"fmt"
	"time"
)

type ServiceOptions struct {
	Port        int
	ConfigFile  string
	GracePeriod time.Duration
}

func (o *ServiceOptions) Validate() error {
	if o.ConfigFile == "" {
		return fmt.Errorf("missing config-file")
	}

	return nil
}

func (o *ServiceOptions) AddFlags(fs *flag.FlagSet) {
	fs.IntVar(&o.Port, "port", 8888, "Port to listen on.")
	fs.StringVar(&o.ConfigFile, "config-file", "", "Path to config file.")
	fs.DurationVar(&o.GracePeriod, "grace-period", 180*time.Second, "On shutdown, try to handle remaining events for the specified duration.")
}
