package ports

// RiskCategory represents a broad classification of a port's role.
type RiskCategory string

const (
	CategorySystem    RiskCategory = "system"
	CategoryDatabase  RiskCategory = "database"
	CategoryWeb       RiskCategory = "web"
	CategoryRemote    RiskCategory = "remote-access"
	CategoryMessaging RiskCategory = "messaging"
	CategoryUnknown   RiskCategory = "unknown"
)

// ClassifyResult holds the classification output for a single port.
type ClassifyResult struct {
	Port     PortInfo
	Category RiskCategory
	Reason   string
}

var categoryRules = []struct {
	ports    []int
	protocol string
	category RiskCategory
	reason   string
}{
	{[]int{22, 23, 3389, 5900}, "tcp", CategoryRemote, "remote access protocol"},
	{[]int{80, 443, 8080, 8443}, "tcp", CategoryWeb, "web server port"},
	{[]int{3306, 5432, 1433, 6379, 27017, 5984}, "tcp", CategoryDatabase, "database service port"},
	{[]int{25, 465, 587, 143, 993, 110, 995}, "tcp", CategoryMessaging, "mail/messaging protocol"},
	{[]int{53, 67, 68, 123, 161, 162}, "", CategorySystem, "core system/network protocol"},
}

// ClassifyPort returns the RiskCategory and a reason string for a given PortInfo.
func ClassifyPort(p PortInfo) ClassifyResult {
	for _, rule := range categoryRules {
		for _, rp := range rule.ports {
			if p.Port == rp && (rule.protocol == "" || rule.protocol == p.Protocol) {
				return ClassifyResult{Port: p, Category: rule.category, Reason: rule.reason}
			}
		}
	}
	return ClassifyResult{Port: p, Category: CategoryUnknown, Reason: "no matching classification rule"}
}

// ClassifyAll classifies a slice of PortInfo values.
func ClassifyAll(ports []PortInfo) []ClassifyResult {
	results := make([]ClassifyResult, 0, len(ports))
	for _, p := range ports {
		results = append(results, ClassifyPort(p))
	}
	return results
}
