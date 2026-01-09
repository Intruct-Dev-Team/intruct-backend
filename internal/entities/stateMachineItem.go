package entities

type StateMachineItem struct {
	ID             int    `db:"item_id"`
	StateMachineID int    `db:"state_machine_id"`
	StateID        int    `db:"state_id"`
	StateName      string `db:"state_name"`
}
