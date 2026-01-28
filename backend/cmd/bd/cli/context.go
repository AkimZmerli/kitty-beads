// Package cli provides shared runtime context for bd CLI subcommands.
// This enables vertical slice architecture by allowing feature packages
// to access the CLI runtime state without importing package main.
package cli

import (
	"context"
	"sync"

	"github.com/spf13/cobra"
	"github.com/steveyegge/beads/internal/rpc"
	"github.com/steveyegge/beads/internal/storage"
)

// Context holds runtime state accessible to all CLI subcommands.
// Populated by main package during initialization.
type Context struct {
	mu sync.RWMutex

	// Core runtime state
	Store        storage.Storage
	DaemonClient *rpc.Client
	RootCtx      context.Context
	Actor        string
	JSONOutput   bool

	// Helper functions (set by main to avoid import cycles)
	ResolvePartialIDFn           func(ctx context.Context, id string) (string, error)
	MarkDirtyFn                  func()
	MarkDirtyAndScheduleFlushFn  func()
	FatalErrorRespectJSONFn      func(format string, args ...interface{})
	OutputJSONFn                 func(v interface{})
	CheckReadonlyFn              func(cmd string)
	EnsureStoreActiveFn          func() error
	FallbackToDirectModeFn       func(reason string) error
	GetActorWithGitFn            func() string
	IssueIDCompletionFn          func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)
}

// global holds the singleton context instance
var global *Context
var once sync.Once

// Get returns the global CLI context.
// Safe to call from any goroutine.
func Get() *Context {
	once.Do(func() {
		global = &Context{}
	})
	return global
}

// SetStore updates the storage backend.
func (c *Context) SetStore(s storage.Storage) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Store = s
}

// GetStore returns the storage backend.
func (c *Context) GetStore() storage.Storage {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Store
}

// SetDaemonClient updates the RPC client.
func (c *Context) SetDaemonClient(client *rpc.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.DaemonClient = client
}

// GetDaemonClient returns the RPC client (nil in direct mode).
func (c *Context) GetDaemonClient() *rpc.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.DaemonClient
}

// SetRootCtx updates the root context.
func (c *Context) SetRootCtx(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.RootCtx = ctx
}

// GetRootCtx returns the root context.
func (c *Context) GetRootCtx() context.Context {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.RootCtx == nil {
		return context.Background()
	}
	return c.RootCtx
}

// SetActor updates the actor name.
func (c *Context) SetActor(actor string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Actor = actor
}

// GetActor returns the actor name.
func (c *Context) GetActor() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Actor
}

// SetJSONOutput updates the JSON output flag.
func (c *Context) SetJSONOutput(json bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.JSONOutput = json
}

// IsJSONOutput returns true if JSON output is enabled.
func (c *Context) IsJSONOutput() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.JSONOutput
}

// ResolvePartialID resolves a partial issue ID to full ID.
func (c *Context) ResolvePartialID(ctx context.Context, id string) (string, error) {
	if c.ResolvePartialIDFn != nil {
		return c.ResolvePartialIDFn(ctx, id)
	}
	return id, nil // fallback: return as-is
}

// MarkDirty marks the database as needing a flush.
func (c *Context) MarkDirty() {
	if c.MarkDirtyFn != nil {
		c.MarkDirtyFn()
	}
}

// MarkDirtyAndScheduleFlush marks the database as dirty and schedules an auto-flush.
func (c *Context) MarkDirtyAndScheduleFlush() {
	if c.MarkDirtyAndScheduleFlushFn != nil {
		c.MarkDirtyAndScheduleFlushFn()
	}
}

// FatalErrorRespectJSON prints error and exits, respecting JSON mode.
func (c *Context) FatalErrorRespectJSON(format string, args ...interface{}) {
	if c.FatalErrorRespectJSONFn != nil {
		c.FatalErrorRespectJSONFn(format, args...)
	}
}

// OutputJSON outputs data as JSON.
func (c *Context) OutputJSON(v interface{}) {
	if c.OutputJSONFn != nil {
		c.OutputJSONFn(v)
	}
}

// CheckReadonly checks if operation is allowed in readonly mode.
func (c *Context) CheckReadonly(cmd string) {
	if c.CheckReadonlyFn != nil {
		c.CheckReadonlyFn(cmd)
	}
}

// EnsureStoreActive checks if the store is ready for use.
func (c *Context) EnsureStoreActive() error {
	if c.EnsureStoreActiveFn != nil {
		return c.EnsureStoreActiveFn()
	}
	return nil
}

// FallbackToDirectMode switches from daemon to direct mode.
func (c *Context) FallbackToDirectMode(reason string) error {
	if c.FallbackToDirectModeFn != nil {
		return c.FallbackToDirectModeFn(reason)
	}
	return nil
}

// GetActorWithGit returns the actor name, preferring git config.
func (c *Context) GetActorWithGit() string {
	if c.GetActorWithGitFn != nil {
		return c.GetActorWithGitFn()
	}
	return c.Actor
}

// IssueIDCompletion provides shell completion for issue IDs.
func (c *Context) IssueIDCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if c.IssueIDCompletionFn != nil {
		return c.IssueIDCompletionFn(cmd, args, toComplete)
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}

// RootCmd is set by main to allow subcommands to register.
var RootCmd *cobra.Command
