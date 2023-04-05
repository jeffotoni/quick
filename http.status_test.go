package quick

import "testing"

func TestQuick_HttpStatus(t *testing.T) {
	t.Log(StatusContinue)
	t.Log(StatusSwitchingProtocols)
	t.Log(StatusProcessing)
	t.Log(StatusEarlyHints)

	t.Log(StatusOK)
	t.Log(StatusCreated)
	t.Log(StatusAccepted)
	t.Log(StatusNonAuthoritativeInfo)
	t.Log(StatusNoContent)
	t.Log(StatusResetContent)
	t.Log(StatusPartialContent)

	t.Log(MethodGet)
	t.Log(MethodHead)
	t.Log(MethodPost)
	t.Log(MethodPut)
	t.Log(MethodPatch)
	t.Log(MethodDelete)
	t.Log(MethodConnect)
	t.Log(MethodTrace)
}
