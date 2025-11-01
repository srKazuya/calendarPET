//Package request provides
package request
import(
	"errors"
)


var (
	ErrInvalidRequest        = errors.New("invalid request")
	ErrEmptyReqBody          = errors.New("request body is empty")
	ErrFailedToDecodeReqBody = errors.New("failed to decode request body")
)
