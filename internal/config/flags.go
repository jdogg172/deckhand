package config

import (
	"fmt"

	"github.com/spf13/pflag"
)

type Flags struct {
	Context    string
	Namespace  string
	Mode       string
	ReadOnly   bool
	Kubeconfig string
	LogLevel   string
	NoColor    bool
	ShowVer    bool
}

func ParseFlags(args []string) (Flags, error) {
	var f Flags

	fs := pflag.NewFlagSet("deckhand", pflag.ContinueOnError)
	fs.StringVar(&f.Context, "context", "", "kube context to use")
	fs.StringVar(&f.Namespace, "namespace", "", "namespace to use")
	fs.StringVar(&f.Mode, "mode", "", "mode to use: ops or pipeline")
	fs.BoolVar(&f.ReadOnly, "readonly", false, "disable mutating actions")
	fs.StringVar(&f.Kubeconfig, "kubeconfig", "", "path to kubeconfig")
	fs.StringVar(&f.LogLevel, "log-level", "info", "log level")
	fs.BoolVar(&f.NoColor, "no-color", false, "disable color")
	fs.BoolVar(&f.ShowVer, "version", false, "show version")

	if err := fs.Parse(args); err != nil {
		return f, err
	}

	if f.Mode != "" && f.Mode != "ops" && f.Mode != "pipeline" {
		return f, fmt.Errorf("invalid mode %q, expected ops or pipeline", f.Mode)
	}

	return f, nil
}
