package event_consumer

import (
	"context"
	"github.com/MangoSociety/know_api/internal/telegram/events"
	"log"
	"sync"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func NewConsumer(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start(ctx context.Context) error {
	var wg sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return ctx.Err()
		default:
			gotEvents, err := c.fetcher.Fetch(ctx, c.batchSize)
			if err != nil {
				log.Printf("[ERR] consumer: %s", err.Error())
				continue
			}

			if len(gotEvents) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}

			wg.Add(1)
			go func(events []events.Event) {
				defer wg.Done()
				if err := c.handleEvents(ctx, events); err != nil {
					log.Print(err)
				}
			}(gotEvents)
		}
	}
}

func (c *Consumer) handleEvents(ctx context.Context, events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(ctx, event); err != nil {
			log.Printf("can't handle event: %s", err.Error())
			continue
		}
	}

	return nil
}
