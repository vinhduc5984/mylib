package utils

import (
	"context"
)

func SendNotify(ctx context.Context, urlOrServiceAddr string) error {
	_, err := RestPostWithContext(urlOrServiceAddr, "/core/notification/v1/send", map[string]interface{}{}, ctx)
	if err != nil {
		return err
	}

	return nil
}
