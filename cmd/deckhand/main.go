package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/jdogg172/deckhand/internal/app"
	"github.com/jdogg172/deckhand/internal/clients"
	"github.com/jdogg172/deckhand/internal/config"
)

func main() {
	ctx := context.Background()

	flags, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse flags: %v\n", err)
		os.Exit(1)
	}

	if flags.ShowVer {
		fmt.Println("deckhand dev")
		return
	}

	cfg, err := config.Load(flags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	kubeFactory, err := clients.NewKubeFactory(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize Kubernetes clients: %v\n", err)
		os.Exit(1)
	}

	if cfg.Context == "" {
		cfg.Context = kubeFactory.CurrentContext
	}
	if cfg.Namespace == "" {
		cfg.Namespace = kubeFactory.CurrentNamespace
	}

	model := app.NewModel(ctx, cfg, kubeFactory)
	prog := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := prog.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "application error: %v\n", err)
		os.Exit(1)
	}
}
