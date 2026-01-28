// Package shared provides shared state and utilities for CLI commands.
// This package enables the vertical slice architecture by allowing command
// packages to access common runtime state without circular dependencies.
package shared

import (
	"context"
	"sync"

	"github.com/steveyegge/beads/internal/hooks"
	"github.com/steveyegge/beads/internal/rpc"
	"github.com/steveyegge/beads/internal/storage"
)

// Runtime holds all shared runtime state for CLI commands.
// Access via the package-level functions (GetStore, GetActor, etc.)
type Runtime struct {
	mu sync.RWMutex

	// Core state
	Store        storage.Storage
	Actor        string
	DBPath       string
	JSONOutput   bool
	DaemonClient *rpc.Client
	HookRunner   *hooks.Runner

	// Context for cancellation
	Ctx    context.Context
	Cancel context.CancelFunc

	// Flags
	NoDaemon     bool
	NoAutoFlush  bool
	NoAutoImport bool
	SandboxMode  bool
	ReadonlyMode bool
	VerboseFlag  bool
	QuietFlag    bool
}

var (
	runtime  *Runtime
	initOnce sync.Once
)

// Init initializes the shared runtime. Called once from main.go.
func Init() *Runtime {
	initOnce.Do(func() {
		runtime = &Runtime{}
	})
	return runtime
}

// Get returns the shared runtime instance.
func Get() *Runtime {
	if runtime == nil {
		panic("shared.Init() must be called before shared.Get()")
	}
	return runtime
}

// GetStore returns the storage instance.
func GetStore() storage.Storage {
	r := Get()
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Store
}

// SetStore sets the storage instance.
func SetStore(s storage.Storage) {
	r := Get()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Store = s
}

// GetActor returns the actor name for audit trails.
func GetActor() string {
	r := Get()
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Actor
}

// SetActor sets the actor name.
func SetActor(actor string) {
	r := Get()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Actor = actor
}

// GetDBPath returns the database path.
func GetDBPath() string {
	r := Get()
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.DBPath
}

// SetDBPath sets the database path.
func SetDBPath(path string) {
	r := Get()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.DBPath = path
}

// IsJSONOutput returns whether JSON output is enabled.
func IsJSONOutput() bool {
	r := Get()
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.JSONOutput
}

// SetJSONOutput sets JSON output mode.
func SetJSONOutput(enabled bool) {
	r := Get()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.JSONOutput = enabled
}

// GetDaemonClient returns the RPC daemon client (may be nil).
func GetDaemonClient() *rpc.Client {
	r := Get()
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.DaemonClient
}

// SetDaemonClient sets the RPC daemon client.
func SetDaemonClient(client *rpc.Client) {
	r := Get()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.DaemonClient = client
}

// GetContext returns the root context for operations.
func GetContext() context.Context {
	r := Get()
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.Ctx == nil {
		return context.Background()
	}
	return r.Ctx
}

// SetContext sets the root context and cancel function.
func SetContext(ctx context.Context, cancel context.CancelFunc) {
	r := Get()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Ctx = ctx
	r.Cancel = cancel
}

// GetHookRunner returns the hook runner.
func GetHookRunner() *hooks.Runner {
	r := Get()
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.HookRunner
}

// SetHookRunner sets the hook runner.
func SetHookRunner(runner *hooks.Runner) {
	r := Get()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.HookRunner = runner
}

// IsReadonly returns whether readonly mode is enabled.
func IsReadonly() bool {
	r := Get()
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.ReadonlyMode
}

// SetReadonly sets readonly mode.
func SetReadonly(readonly bool) {
	r := Get()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ReadonlyMode = readonly
}

// IsVerbose returns whether verbose output is enabled.
func IsVerbose() bool {
	r := Get()
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.VerboseFlag
}

// IsQuiet returns whether quiet mode is enabled.
func IsQuiet() bool {
	r := Get()
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.QuietFlag
}
