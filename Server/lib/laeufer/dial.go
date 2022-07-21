package laeufer

import (
	"context"
	"net"
	"time"

	"github.com/c-f/talje/lib/utils"
)

type DialFunc func(network, address string) (net.Conn, error)
type DialContextFunc func(ctx context.Context, network, address string) (net.Conn, error)

// conn
func DialWithRetry(dialctx DialContextFunc) DialFunc {
	return func(network, address string) (net.Conn, error) {
		var newConn net.Conn

		retryErr := utils.Retry(3, 10*time.Second, func() (err error) {

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			newConn, err = dialctx(ctx, network, address)
			return err
		})
		if retryErr != nil {
			return nil, retryErr
		}

		return newConn, nil

	}
}
