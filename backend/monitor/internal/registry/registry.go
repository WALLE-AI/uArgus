package registry

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
)

// Registry is the central container that holds all Sources,
// manages their scheduling and tracks their health.
type Registry struct {
	sources   map[string]Source
	order     []string // topologically sorted names
	health    *HealthTracker
	scheduler *Scheduler
}

// New creates an empty Registry.
func New() *Registry {
	r := &Registry{
		sources: make(map[string]Source),
		health:  NewHealthTracker(),
	}
	r.scheduler = NewScheduler(r.dispatchRun)
	return r
}

// Register adds a Source. Panics on duplicate name.
func (r *Registry) Register(s Source) {
	name := s.Name()
	if _, dup := r.sources[name]; dup {
		panic(fmt.Sprintf("registry: duplicate source name %q", name))
	}
	r.sources[name] = s
	r.health.Register(s)
}

// Boot performs topological sort on dependencies, registers schedules, and starts.
func (r *Registry) Boot(ctx context.Context) error {
	order, err := r.topoSort()
	if err != nil {
		return fmt.Errorf("registry boot: %w", err)
	}
	r.order = order

	for _, name := range r.order {
		s := r.sources[name]
		if err := r.scheduler.Add(s); err != nil {
			return fmt.Errorf("registry: schedule %s: %w", name, err)
		}
	}

	r.scheduler.Start(ctx)
	slog.Info("registry booted", "sources", len(r.sources), "order", r.order)
	return nil
}

// Shutdown gracefully stops all scheduling.
func (r *Registry) Shutdown(_ context.Context) error {
	r.scheduler.Stop()
	slog.Info("registry shutdown")
	return nil
}

// HealthSnapshot returns a copy of all source health statuses.
func (r *Registry) HealthSnapshot() map[string]HealthStatus {
	return r.health.Snapshot()
}

// TriggerOnDemand sends a trigger to an OnDemand source.
func (r *Registry) TriggerOnDemand(name string) error {
	if !r.scheduler.Trigger(name) {
		return fmt.Errorf("registry: source %q not found or not on-demand", name)
	}
	return nil
}

// dispatchRun is the callback invoked by Scheduler on each trigger.
func (r *Registry) dispatchRun(ctx context.Context, s Source) {
	name := s.Name()
	r.health.RecordAttempt(name)

	result, err := s.Run(ctx)
	if err != nil {
		slog.Error("source run failed", "source", name, "err", err)
		r.health.RecordFailure(name)
		return
	}

	r.health.RecordSuccess(name, result.Metrics)
	slog.Info("source run succeeded", "source", name,
		"records", result.Metrics.RecordCount,
		"duration", result.Metrics.Duration)
}

// ── topological sort ────────────────────────────────────────

func (r *Registry) topoSort() ([]string, error) {
	inDeg := make(map[string]int)
	adj := make(map[string][]string)

	for name := range r.sources {
		if _, ok := inDeg[name]; !ok {
			inDeg[name] = 0
		}
		for _, dep := range r.sources[name].Dependencies() {
			adj[dep] = append(adj[dep], name)
			inDeg[name]++
		}
	}

	var queue []string
	for name, d := range inDeg {
		if d == 0 {
			queue = append(queue, name)
		}
	}
	sort.Strings(queue) // deterministic order for same-level

	var result []string
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		result = append(result, n)
		for _, next := range adj[n] {
			inDeg[next]--
			if inDeg[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	if len(result) != len(r.sources) {
		return nil, fmt.Errorf("dependency cycle detected")
	}
	return result, nil
}
