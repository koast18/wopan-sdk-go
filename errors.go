package wopan

import "fmt"

// errorFromStatus and errorFromRsp keep the error strings of the dispatcher in
// one place, so the plain (non-dispatcher) endpoints can report the same shape.
func errorFromStatus(op, status string) error {
	return fmt.Errorf("%s failed with status: %s", op, status)
}

func errorFromRsp(op, code, desc string) error {
	return fmt.Errorf("%s failed with rsp_code: %s, rep_desc: %s", op, code, desc)
}
