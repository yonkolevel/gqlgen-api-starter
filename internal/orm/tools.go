package orm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/txbrown/gqlgen-api-starter/internal/gql/model"
	"github.com/txbrown/gqlgen-api-starter/pkg/utils"
	"gorm.io/gorm"
)

// ParseFilters parses the filter and adds the where condition to the transaction
func ParseFilters(db *gorm.DB, filters []*model.QueryFilter) (*gorm.DB, error) {
	for _, f := range filters {
		condition := utils.ToSnakeCase(f.Field) + " " + opToSQL(f.Op)
		switch f.Op {
		case model.OperationTypeBetween:
			if len(f.Values) != 2 {
				return db, errors.New("Operation [" + f.Op.String() +
					"] needs an array with exactly two items in [values] field")
			}
			if f.LinkOperation != nil && *f.LinkOperation == model.LinkOperationTypeOr {
				db = db.Or(condition, f.Values[0], f.Values[1])
			} else {
				db = db.Where(condition, f.Values[0], f.Values[1])
			}
		case model.OperationTypeIn, model.OperationTypeNotIn:
			if len(f.Values) < 1 {
				return db, errors.New("Operation [" + f.Op.String() +
					"] needs an array with at least 1 item on [values] field")
			}
			if f.LinkOperation != nil && *f.LinkOperation == model.LinkOperationTypeOr {
				db = db.Or(condition, f.Values)
			} else {
				db = db.Where(condition, f.Values)
			}
		case model.OperationTypeMatch:
			if f.LinkOperation != nil && *f.LinkOperation == model.LinkOperationTypeOr {
				db = db.Or("MATCH("+utils.ToSnakeCase(f.Field)+
					") AGAINST (? IN BOOLEAN MODE)", f.Value)
			} else {
				db = db.Where("MATCH("+utils.ToSnakeCase(f.Field)+
					") AGAINST (? IN BOOLEAN MODE)", f.Value)
			}
		case model.OperationTypeIsNotNull:
			fallthrough
		case model.OperationTypeIsNull:
			if f.LinkOperation != nil && *f.LinkOperation == model.LinkOperationTypeOr {
				db = db.Or(condition)
			} else {
				db = db.Where(condition)
			}

		default:
			if f.Value == nil {
				return db, errors.New("Operation [" + f.Op.String() +
					"] needs the field [value] to compare")
			}
			if f.LinkOperation != nil && *f.LinkOperation == model.LinkOperationTypeOr {
				db = db.Or(condition, f.Value)
			} else {
				db = db.Where(condition, f.Value)
			}
		}
	}
	return db, db.Error
}

func opToSQL(op model.OperationType) string {
	return map[model.OperationType]string{
		model.OperationTypeEquals:           " = ?",
		model.OperationTypeNotEquals:        " != ?",
		model.OperationTypeLessThan:         " < ?",
		model.OperationTypeLessThanEqual:    " <= ?",
		model.OperationTypeGreaterThan:      " > ?",
		model.OperationTypeGreaterThanEqual: " >= ?",
		model.OperationTypeIs:               " IS ?",
		model.OperationTypeIsNull:           " IS NULL",
		model.OperationTypeIsNotNull:        " IS NOT NULL",
		model.OperationTypeIn:               " IN (?)",
		model.OperationTypeNotIn:            " NOT IN (?)",
		model.OperationTypeLike:             " LIKE ?",
		model.OperationTypeILike:            " ILIKE ?",
		model.OperationTypeNotLike:          " NOT LIKE ?",
		model.OperationTypeBetween:          " BETWEEN ? AND ?",
	}[op]
}

func arrayToString(a interface{}, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), "' '", delim, -1), "[]")
}
