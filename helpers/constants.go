package helpers

const (
	CreatedMessage = "Successfully Inserted!"
)

type ContextValues struct {
	UserID     string
	AccessUUID string
	Role       string
}

type EnviormentVariables struct {
	MONGO_USERNAME     string
	MONGO_PASSWORD     string
	MONGO_DBNAME       string
	SENDGRID_API_KEY   string
	JWT_ACCESS_SECRET  string
	JWT_REFRESH_SECRET string
	REDIS_DSN          string
	REDIS_PASSWORD     string
}

var CtxValues = ContextValues{
	UserID:     "userID",
	AccessUUID: "accessUUID",
	Role:       "role",
}

var EnvKeys = EnviormentVariables{
	MONGO_USERNAME:     "MONGO_USERNAME",
	MONGO_PASSWORD:     "MONGO_PASSWORD",
	MONGO_DBNAME:       "MONGO_DBNAME",
	SENDGRID_API_KEY:   "SENDGRID_API_KEY",
	JWT_ACCESS_SECRET:  "JWT_ACCESS_SECRET",
	JWT_REFRESH_SECRET: "JWT_REFRESH_SECRET",
	REDIS_DSN:          "REDIS_DSN",
	REDIS_PASSWORD:     "REDIS_PASSWORD",
}
