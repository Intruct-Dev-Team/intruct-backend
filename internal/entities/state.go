package entities

type StateName string

const (
	InCreation StateName = "in creation"
	Failed     StateName = "failed"
	Created    StateName = "created"
	Published  StateName = "published"

// Pending    StateName = "pending"
// Approved   StateName = "approved"
// Rejected   StateName = "rejected"
// Planned    StateName = "planned"
// Cancelled  StateName = "cancelled"
// Ongoing    StateName = "ongoing"
// Finished   StateName = "finished"
// Conflicted StateName = "conflicted"
// Completed  StateName = "completed"
)

type State struct {
	ID   int       `db:"state_id"`
	Name StateName `db:"name"`
}
