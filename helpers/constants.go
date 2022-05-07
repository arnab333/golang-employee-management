package helpers

const (
	CreatedMessage = "Successfully Inserted!"
)

type userRoles struct {
	Superadmin string `json:"superadmin"`
	Admin      string `json:"admin"`
	User       string `json:"user"`
}

var UserRoles userRoles = userRoles{
	Superadmin: "superadmin",
	Admin:      "admin",
	User:       "user",
}
