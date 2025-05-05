package logger

type Category string
type SubCategory string
type ExtraCategory string

const (
	Server   Category = "Server"
	Postgres Category = "Postgres"
	Redis    Category = "Redis"
	User     Category = "User"
	Admin    Category = "Admin"
	Tweet    Category = "Tweet"
	Config   Category = "Config"
)

const (
	Add     SubCategory = "Add"
	Start   SubCategory = "Start"
	Get     SubCategory = "Get"
	Connect SubCategory = "Connect"
	Delete  SubCategory = "Delete"
	New     SubCategory = "New"
	Update  SubCategory = "Update"
	Enter   SubCategory = "Enter"
	Edit    SubCategory = "Edit"
)