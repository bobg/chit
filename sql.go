package chit

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

// QueryerContext is a minimal interface satisfied by *sql.DB
// (from database/sql).
type QueryerContext interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
}

// SQL[T] performs a query against db and returns the results as an iterator of type T.
// T must be a struct type whose fields have the same types,
// in the same order,
// as the values being queried.
// The values produced by the iterator will be instances of that struct type,
// with fields populated by the queried values.
func SQL[T any](ctx context.Context, db QueryerContext, query string, args ...any) *Iter[T] {
	return New(ctx, func(ctx context.Context, ch chan<- T) error {
		var t T
		tt := reflect.TypeOf(t)
		if tt.Kind() != reflect.Struct {
			return fmt.Errorf("type parameter to SQL has %s kind but must be struct", tt.Kind())
		}
		nfields := tt.NumField()

		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var (
				row    T
				rowval = reflect.ValueOf(row)
				ptrs   = make([]interface{}, 0, nfields)
			)
			for i := 0; i < nfields; i++ {
				addr := rowval.Field(i).Addr()
				ptrs = append(ptrs, addr.Interface())
			}
			err = rows.Scan(ptrs...)
			if err != nil {
				return err
			}
			err = chwrite(ctx, ch, row)
			if err != nil {
				return err
			}
		}
		return rows.Err()
	})
}
