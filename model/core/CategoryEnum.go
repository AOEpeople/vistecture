package core

type Category int

const (
	CORE		Category = iota
	PROJECT
	INDIVIDUAL
	EXTERNAL
)

func (category Category) Value() string {
	names := [...]string{
		"core",
		"project",
		"individual",
		"external"}
	if category < CORE || category > EXTERNAL {
		return "Unknown"
	}
	return names[category]
}
