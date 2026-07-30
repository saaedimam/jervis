// Package calendar provides Calendar API integration.
//
// Purpose: Synchronize local platform state with external calendar providers (iCal, Google Calendar).
// Responsibilities:
// - Fetching events from remote providers.
// - Pushing platform-generated events to remote providers.
// - Reconciling conflicts between local and remote states.
// Dependencies: internal/runtime/eventbus/contracts, net/http
// Layer ownership: Service Layer (Layer 3)
package calendar
