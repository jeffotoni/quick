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
}
