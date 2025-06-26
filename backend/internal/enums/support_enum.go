package enums 

type SupportStatus string

const (
	OpenTicket     SupportStatus = "open"
	CloseTicket    SupportStatus = "closed"
	ResolvedTicket SupportStatus = "resolved"
	PendingTicket  SupportStatus = "pending"
)

type Priority string

const (
	Low    Priority = "low"
	Medium Priority = "medium"
	High   Priority = "high"
)