package consts

import (
	"fmt"
	"strings"

	"github.com/txbrown/gqlgen-api-starter/pkg/utils"
)

type permissionTypes struct {
	Create string
	Read   string
	Update string
	Delete string
	List   string
	Assign string
	Upload string
}

type entitynames struct {
	Users           string
	Roles           string
	Permissions     string
	RoleParents     string
	RolePermissions string
	UserPermissions string
	UserProfiles    string
	UserRoles       string
	Pack            string
	Circuit         string
	Lesson          string
	File            string
	PackDetails     string
	PackReview      string
	PackActivity    string
	Category        string
	Feedback        string
	Trophy          string
	UserAPIKeys     string
}

type role struct {
	Name        string
	Description string
}

type dialects struct {
	PostgresSQL string
	MySQL       string
}

type authProviders struct {
	DB string
}

var (
	// Permissions has the types of permissions that can be assigned
	Permissions = permissionTypes{
		Create: "create:%s",
		Read:   "read:%s",
		Update: "update:%s",
		Delete: "delete:%s",
		List:   "list:%s",
		Assign: "assign:%s",
		Upload: "upload:%s",
	}
	// EntityNames the names of the tables in the server
	EntityNames = entitynames{
		Users:           "Users",
		Roles:           "Roles",
		Permissions:     "Permissions",
		RoleParents:     "RoleParents",
		RolePermissions: "RolePermissions",
		UserPermissions: "UserPermissions",
		UserProfiles:    "UserProfiles",
		UserRoles:       "UserRoles",
		Pack:            "Packs",
		Circuit:         "Circuits",
		Lesson:          "Lessons",
		File:            "Files",
		PackDetails:     "PackDetails",
		Category:        "Categories",
		PackReview:      "PackReviews",
		PackActivity:    "PackActivities",
		Feedback:        "Feedbacks",
		Trophy:          "Trophies",
		UserAPIKeys:     "UserAPIKeys",
	}
	// Dialects are definition of databases
	Dialects = dialects{
		PostgresSQL: "postgres",
		MySQL:       "mysql",
	}

	// Roles that are part of the systme
	Roles = []role{
		{
			Name:        "admin",
			Description: "Administrator of the app",
		},
		{
			Name:        "user",
			Description: "Normal user of the app",
		},
	}

	// Providers
	Providers = authProviders{
		DB: "db",
	}

	NestedFmt = "%s.%s"
)

// GetTableName gets the db normalized tablename
func GetTableName(tablename string) string {
	return utils.ToSnakeCase(tablename)
}

// FormatPermissionTag returns a string formatted action:entity permission
func FormatPermissionTag(action string, entity string) string {
	return fmt.Sprintf(action, entity)
}

// FormatPermissionDesc returns a string with the description of the
// action:entity permission
func FormatPermissionDesc(action string, entity string) string {
	return "Allows the user to " +
		strings.ReplaceAll(FormatPermissionTag(action, entity), ":", " ")
}
