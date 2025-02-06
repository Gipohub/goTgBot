package e

import "fmt"

//унифицируем вариотизацию возвращаемых ошибок
func Wrap(msg string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}
