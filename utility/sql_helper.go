package utility

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
)

func GetStringValue(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

func GetIntValue(s sql.NullInt64) int64 {
	if s.Valid {
		return s.Int64
	}
	return 0
}

func GetBoolValue(s sql.NullBool) bool {
	if s.Valid {
		return s.Bool
	}
	return false
}

func GetFloatValue(s sql.NullFloat64) float64 {
	if s.Valid {
		return s.Float64
	}
	return 0.0
}

func Transact(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(tx)
	return err
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 || s == "0000-00-00" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewNullInt(i int64) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

func NewNullFloat(i float64) sql.NullFloat64 {
	if i == 0 {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{
		Float64: i,
		Valid:   true,
	}
}

func MapSqlValues(s interface{}, t interface{}) {
	source := reflect.ValueOf(s).Elem()
	target := reflect.ValueOf(t).Elem()
	for j := 0; j < source.NumField(); j++ {
		field := source.Field(j)
		name := source.Type().Field(j).Name
		fi := field.Interface()
		val := reflect.ValueOf(fi)
		for y := 0; y < target.NumField(); y++ {
			tfield := target.Field(y)
			tname := target.Type().Field(y).Name
			if name == tname {
				switch field.Type().String() {
				case "sql.NullString":
					if tfield.CanSet() {
						tfield.SetString(GetStringValue(val.Interface().(sql.NullString)))
					}
				case "sql.NullInt64":
					if tfield.CanSet() {
						tfield.SetInt(GetIntValue(val.Interface().(sql.NullInt64)))
					}
				case "sql.NullFloat64":
					if tfield.CanSet() {
						tfield.SetFloat(GetFloatValue(val.Interface().(sql.NullFloat64)))
					}
				case "sql.NullBool":
					if tfield.CanSet() {
						tfield.SetBool(GetBoolValue(val.Interface().(sql.NullBool)))
					}
				}
				break
			}
		}
	}
}
