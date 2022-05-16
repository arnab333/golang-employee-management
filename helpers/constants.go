package helpers

type ContextValues struct {
	UserID      string
	AccessUUID  string
	Role        string
	Permissions string
}

type EnviormentVariables struct {
	MONGO_USERNAME      string
	MONGO_PASSWORD      string
	MONGO_DBNAME        string
	SENDGRID_API_KEY    string
	JWT_ACCESS_SECRET   string
	JWT_REFRESH_SECRET  string
	REDIS_DSN           string
	REDIS_PASSWORD      string
	APP_ENV             string
	GOOGLE_API_KEY      string
	GOOGLE_CALENDAR_ID  string
	SENDGRID_FROM_EMAIL string
}

type RolesList struct {
	SuperAdmin        string
	Admin             string
	Accountant        string
	AccountantManager string
	Manager           string
	HR                string
	HRManager         string
}

type PermissionsList struct {
	CreateRole string
	ReadRole   string
	UpdateRole string
	DeleteRole string

	CreateHoliday string
	ReadHoliday   string
	UpdateHoliday string
	DeleteHoliday string

	CreateUser string
	ReadUser   string
	UpdateUser string
	DeleteUser string
}

var CtxValues = ContextValues{
	UserID:      "userID",
	AccessUUID:  "accessUUID",
	Role:        "role",
	Permissions: "permissions",
}

var EnvKeys = EnviormentVariables{
	MONGO_USERNAME:      "MONGO_USERNAME",
	MONGO_PASSWORD:      "MONGO_PASSWORD",
	MONGO_DBNAME:        "MONGO_DBNAME",
	SENDGRID_API_KEY:    "SENDGRID_API_KEY",
	JWT_ACCESS_SECRET:   "JWT_ACCESS_SECRET",
	JWT_REFRESH_SECRET:  "JWT_REFRESH_SECRET",
	REDIS_DSN:           "REDIS_DSN",
	REDIS_PASSWORD:      "REDIS_PASSWORD",
	APP_ENV:             "APP_ENV",
	GOOGLE_API_KEY:      "GOOGLE_API_KEY",
	GOOGLE_CALENDAR_ID:  "GOOGLE_CALENDAR_ID",
	SENDGRID_FROM_EMAIL: "SENDGRID_FROM_EMAIL",
}

var UserRoles = RolesList{
	SuperAdmin:        "superAdmin",
	Admin:             "admin",
	Accountant:        "accountant",
	AccountantManager: "accountantManager",
	Manager:           "manager",
	HR:                "hr",
	HRManager:         "hrManager",
}

var UserPermissions = PermissionsList{
	CreateRole: "CreateRole",
	ReadRole:   "ReadRole",
	UpdateRole: "UpdateRole",
	DeleteRole: "DeleteRole",

	CreateHoliday: "CreateHoliday",
	ReadHoliday:   "ReadHoliday",
	UpdateHoliday: "UpdateHoliday",
	DeleteHoliday: "DeleteHoliday",

	CreateUser: "CreateUser",
	ReadUser:   "ReadUser",
	UpdateUser: "UpdateUser",
	DeleteUser: "DeleteUser",
}

const (
	CreatedMessage = "Successfully Inserted!"
	Unauthorized   = "You are not authorized!"
	InvalidID      = "Invalid ID!"
	RequiredID     = "ID is required!"
	RequiredField  = "Field is required!"
)
