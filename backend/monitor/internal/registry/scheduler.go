package registry

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// Scheduler manages Cron, Interval and OnDemand dispatch for registered Sources.
type Scheduler struct {
	mu       sync.Mutex
	cron     *cron.Cron
	tickers  []*tickerEntry
	onDemand map[string]chan struct{}
	dispatch func(ctx context.Context, s Source)
}

type tickerEntry struct {
	source Source
	every  time.Duration
	cancel context.CancelFunc
}

// NewScheduler creates a Scheduler. dispatch is called on each trigger.
func NewScheduler(dispatch func(ctx context.Context, s Source)) *Scheduler {
	return &Scheduler{
		cron:     cron.New(cron.WithSeconds()),
		onDemand: make(map[string]chan struct{}),
		dispatch: dispatch,
	}
}

// Add registers a Source according to its Schedule type.
func (sc *Scheduler) Add(s Source) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	switch sched := s.Spec().Schedule.(type) {
	case CronSchedule:
		_, err := sc.cron.AddFunc(sched.Expr, func() {
			sc.dispatch(context.Background(), s)
		})
		if err != nil {
			return err
		}
	case IntervalSchedule:
		sc.tickers = append(sc.tickers, &tickerEntry{source: s, every: sched.Every})
	case OnDemandSchedule:
		ch := make(chan struct{}, 1)
		sc.onDemand[s.Name()] = ch
	default:
		slog.Warn("unknown schedule type", "source", s.Name())
	}
	return nil
}

// Start begins all cron jobs and interval tickers.
func (sc *Scheduler) Start(ctx context.Context) {
	sc.cron.Start()

	for _, te := range sc.tickers {
		te := te
		tickCtx, cancel := context.WithCancel(ctx)
		te.cancel = cancel
		go func() {
			ticker := time.NewTicker(te.every)
			defer ticker.Stop()
			// fire immediately on start
			sc.dispatch(tickCtx, te.source)
			for {
				select {
				case <-tickCtx.Done():
					return
				case <-ticker.C:
					sc.dispatch(tickCtx, te.source)
				}
			}
		}()
	}

	for name, ch := range sc.onDemand {
		name, ch := name, ch
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case <-ch:
					sc.mu.Lock()
					// find source by iterating — small N so OK
					sc.mu.Unlock()
					slog.Info("on-demand triggered", "source", name)
				}
			}
		}()
	}
}

// Trigger sends a signal to an OnDemand source.
func (sc *Scheduler) Trigger(name string) bool {
	sc.mu.Lock()
	ch, ok := sc.onDemand[name]
	sc.mu.Unlock()
	if !ok {
		return false
	}
	select {
	case ch <- struct{}{}:
	default:
	}
	return true
}

// Stop halts all scheduling.
func (sc *Scheduler) Stop() {
	sc.cron.Stop()
	for _, te := range sc.tickers {
		if te.cancel != nil {
			te.cancel()
		}
	}
}
