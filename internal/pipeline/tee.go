/*
Package pipeline implements the Tee concurrency pattern for QueryEvent streams.
It duplicates each event onto two independent output channels so downstream
consumers can process the same input in parallel without interfering
with each other.
*/
package pipeline

import (
	"context"

	"github.com/atharvamhaske/supaspy/internal/models"
)

// Done wraps a channel so reads stop cleanly when the context is cancelled
// it prevents a goroutine leaks when the pipeline shuts down.
func Done(ctx context.Context, in <-chan models.QueryEvent) <-chan models.QueryEvent {
	out := make(chan models.QueryEvent)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}
