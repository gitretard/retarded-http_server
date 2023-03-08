package std

import (
	"errors"
	"strings"
)
func (h* Req) ParseFormData() error{
	t1 := strings.Split(h.Data.Body,"&")
	for _, i := range t1 {
        t2 := strings.Split(i, "=")
        if len(t2) != 2 {
            return errors.New("Malformed Data!")
        }
        key, err := URLunescape(t2[0])
        if err != nil {
            return err
        }
        val, err := URLunescape(t2[1])
        if err != nil {
            return err
        }
        h.Data.FormData[key] = val
    }
	return nil
}