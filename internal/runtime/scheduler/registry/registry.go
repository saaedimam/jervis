package registry

import (
	"sync"
	"github.com/ioriimasu/jervis/internal/runtime/scheduler/contracts"
	"github.com/ioriimasu/jervis/internal/runtime/scheduler/errors"
)

type Registry struct {
	mu   sync.RWMutex
	jobs []contracts.Job
	ids  map[string]int // maps ID to index in jobs slice
}

func New() contracts.Registry {
	return &Registry{
		jobs: make([]contracts.Job, 0),
		ids:  make(map[string]int),
	}
}

func (r *Registry) Register(job contracts.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := job.ID()
	if _, exists := r.ids[id]; exists {
		return errors.ErrJobAlreadyExists
	}

	r.ids[id] = len(r.jobs)
	r.jobs = append(r.jobs, job)
	return nil
}

func (r *Registry) Unregister(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	idx, exists := r.ids[id]
	if !exists {
		return errors.ErrJobNotFound
	}

	// Remove from map
	delete(r.ids, id)

	// Remove from slice
	r.jobs = append(r.jobs[:idx], r.jobs[idx+1:]...)

	// Update indices for subsequent jobs
	for i := idx; i < len(r.jobs); i++ {
		r.ids[r.jobs[i].ID()] = i
	}

	return nil
}

func (r *Registry) Get(id string) (contracts.Job, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	idx, exists := r.ids[id]
	if !exists {
		return nil, false
	}
	return r.jobs[idx], true
}

func (r *Registry) All() []contracts.Job {
	r.mu.RLock()
	defer r.mu.RUnlock()

	jobs := make([]contracts.Job, len(r.jobs))
	copy(jobs, r.jobs)
	return jobs
}

func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.jobs)
}

func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs = make([]contracts.Job, 0)
	r.ids = make(map[string]int)
}
