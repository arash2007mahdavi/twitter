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
	Comment  Category = "Comment"
	Config   Category = "Config"
	Otp      Category = "Otp"
	Error    Category = "Error"
)

const (
	SimpleError     SubCategory = "SimpleError"
	ValidationError SubCategory = "ValidationError"
	Add             SubCategory = "Add"
	Start           SubCategory = "Start"
	Get             SubCategory = "Get"
	Connect         SubCategory = "Connect"
	Delete          SubCategory = "Delete"
	New             SubCategory = "New"
	Update          SubCategory = "Update"
	Enter           SubCategory = "Enter"
	Edit            SubCategory = "Edit"
	See             SubCategory = "See"
	Set             SubCategory = "Set"
	Validate        SubCategory = "Validate"
	Follow          SubCategory = "Follow"
	UnFollow        SubCategory = "UnFollow"
	Follower        SubCategory = "Follower"
	Following       SubCategory = "Following"
	Profile         SubCategory = "Profile"
	Like            SubCategory = "Like"
	Dislike         SubCategory = "Dislike"
)

const (
	Username     ExtraCategory = "Username"
	Userid       ExtraCategory = "Userid"
	Targetid     ExtraCategory = "Targetid"
	Tweetid      ExtraCategory = "Tweetid"
	Commentid    ExtraCategory = "Commentid"
	StatusCode   ExtraCategory = "StatusCode"
	Err          ExtraCategory = "Err"
	MobileNumber ExtraCategory = "MobileNumber"
	Adminname    ExtraCategory = "Admin"
)
