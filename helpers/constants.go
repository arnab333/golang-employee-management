package helpers

const (
	CreatedMessage = "Successfully Inserted!"
)

type ContextValues struct {
	UserID     string
	AccessUUID string
}

var CtxValues = ContextValues{
	UserID:     "userID",
	AccessUUID: "accessUUID",
}
