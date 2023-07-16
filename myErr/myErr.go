package myErr

import "fmt"

func Handle(msg string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", msg, err)
}
