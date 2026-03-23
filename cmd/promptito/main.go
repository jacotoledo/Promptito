package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jtg365/promptito/internal/server"
	"github.com/jtg365/promptito/internal/storage"
)

const (
	version = "2.0.0"
	name    = "promptito"
	desc    = "PROMPT + RAPIDITO - The fastest way to manage AI prompts"
	author  = "Jaco Toledo"
	company = "JTG 365 LLC"
	website = "https://jtg365.com"
	github  = "https://github.com/jacotoledo"
)

type Config struct {
	Address      string
	PromptDir    string
	StaticDir    string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type Option func(*Config)

func WithAddress(addr string) Option {
	return func(c *Config) {
		c.Address = addr
	}
}

func WithPromptDir(dir string) Option {
	return func(c *Config) {
		c.PromptDir = dir
	}
}

func WithStaticDir(dir string) Option {
	return func(c *Config) {
		c.StaticDir = dir
	}
}

func WithTimeouts(read, write, idle time.Duration) Option {
	return func(c *Config) {
		c.ReadTimeout = read
		c.WriteTimeout = write
		c.IdleTimeout = idle
	}
}

func NewConfig(opts ...Option) *Config {
	cfg := &Config{
		Address:      ":80",
		PromptDir:    "public/prompts",
		StaticDir:    "public",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func (c *Config) ApplyFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.Address, "addr", c.Address, "Address to listen on")
	fs.StringVar(&c.PromptDir, "prompts", c.PromptDir, "Directory containing skill prompts")
	fs.StringVar(&c.StaticDir, "static", c.StaticDir, "Directory containing static files")
}

func run(ctx context.Context, cfg *Config) error {
	log.Printf("Starting %s v%s", name, version)

	store, err := storage.New(storage.Config{
		Directory: cfg.PromptDir,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	srv, err := server.New(
		server.WithStorage(store),
		server.WithStatic(cfg.StaticDir),
		server.WithPromptDir(cfg.PromptDir),
	)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	httpSrv := &http.Server{
		Addr:         cfg.Address,
		Handler:      srv,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		ErrorLog:     log.New(os.Stderr, "http: ", log.LstdFlags),
	}

	go func() {
		log.Printf("Server listening on %s", cfg.Address)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	log.Println("Server stopped gracefully")
	return nil
}

func main() {
	cfg := NewConfig()

	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\nOptions:\n", os.Args[0])
		fs.PrintDefaults()
		fs.PrintDefaults()
	}

	versionFlag := fs.Bool("version", false, "Print version information")
	helpFlag := fs.Bool("h", false, "Show this help message")
	fs.Bool("help", false, "Show this help message")

	cfg.ApplyFlags(fs)

	if err := fs.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		log.Fatalf("flag error: %v", err)
	}

	if *versionFlag {
		fmt.Printf("%s v%s\n%s\nMade by %s (%s)\n", name, version, desc, author, website)
		os.Exit(0)
	}

	if *helpFlag {
		fs.Usage()
		os.Exit(0)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v", sig)
		cancel()
	}()

	if err := run(ctx, cfg); err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}
