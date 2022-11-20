package recovermw

import (
	"net/http"

	"go.uber.org/zap"
)

type recoverMW struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) *recoverMW {
	return &recoverMW{
		logger: logger,
	}
}

func (r *recoverMW) Middleware(next http.Handler) http.Handler {
	f := func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}

			r.logger.Error("recovered", zap.Any("error", err))
		}()

		next.ServeHTTP(writer, request)
	}

	return http.HandlerFunc(f)
}
