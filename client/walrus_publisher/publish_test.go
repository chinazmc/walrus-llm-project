package walrus_publisher

import (
	"context"
	"testing"
)

func TestPublishWalrus(t *testing.T) {
	PublishWalrus(context.Background(), "hello walrus!!!")
}
