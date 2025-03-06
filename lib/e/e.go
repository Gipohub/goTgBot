package e

import "fmt"

//унифицируем вариотизацию возвращаемых ошибок
func Wrap(msg string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}
func WrapNew(msg string) error {

	return fmt.Errorf("handmade err: %s", msg)
}
