package registry

/* how service registration looks like*/

type Registration struct {
	ServiceName ServiceName
	ServiceURL  string
	//list of services required by this service
	RequiredServices []ServiceName
	//for registry service to talk back to the client service i.e to receive updates
	ServiceUpdateURL string
}

type patchEntry struct {
	Name ServiceName
	URL  string
}

//to communicate both addition and subtraction in single call
type patch struct {
	//services coming online
	Added []patchEntry
	//services going offline
	Removed []patchEntry
}

type ServiceName string

const (
	LogService     = ServiceName("LogService")
	GradingService = ServiceName("GradingService")
)
