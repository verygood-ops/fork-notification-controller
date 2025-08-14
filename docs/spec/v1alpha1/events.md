# Event

The `Event` API defines what information a report of an event issued by a controller should contain.

## Specification

Spec:

```go
type Event struct {
	// The object that this event is about.
	// +required
	InvolvedObject corev1.ObjectReference `json:"involvedObject"`

	// Severity type of this event (info, error)
	// +required
	Severity string `json:"severity"`

	// The time at which this event was recorded.
	// +required
	Timestamp metav1.Time `json:"timestamp"`

	// A human-readable description of this event.
	// Maximum length 39,000 characters
	// +required
	Message string `json:"message"`

	// A machine understandable string that gives the reason
	// for the transition into the object's current status.
	// +required
	Reason string `json:"reason"`

	// Metadata of this event, e.g. apply change set.
	// +optional
	Metadata map[string]string `json:"metadata,omitempty"`

	// Name of the controller that emitted this event, e.g. `source-controller`.
	// +required
	ReportingController string `json:"reportingController"`

	// ID of the controller instance, e.g. `source-controller-xyzf`.
	// +optional
	ReportingInstance string `json:"reportingInstance,omitempty"`
}
```

Event severity:

```go
const (
	EventSeverityInfo string = "info"
	EventSeverityError string = "error"
)
```

Controller implementations can use the [fluxcd/pkg/recorder](https://github.com/fluxcd/pkg/tree/master/recorder)
package to push events to notification-controller API.
