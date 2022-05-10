package helpers

const (
	CreatedMessage = "Successfully Inserted!"
)

type ContextValues struct {
	UserID     string
	AccessUUID string
	Role       string
}

var CtxValues = ContextValues{
	UserID:     "userID",
	AccessUUID: "accessUUID",
	Role:       "role",
}
