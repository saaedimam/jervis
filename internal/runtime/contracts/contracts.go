package contracts

// Validator defines objects that can validate their state.
type Validator interface {
	Validate() error
}

// Freezable defines objects that can be frozen into an immutable state.
type Freezable interface {
	IsFrozen() bool
}

// VersionProvider defines components capable of supplying version information.
type VersionProvider interface {
	Version() string
	GitCommit() string
	BuildDate() string
}

// LifecycleManager defines the contract for deterministic state machine managers.
type LifecycleManager interface {
	State() string
	Start() error
	Stop() error
	TransitionTo(targetState string) error
	Fail(err error) error
}
