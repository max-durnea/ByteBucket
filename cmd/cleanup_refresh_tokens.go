package main

import (
	"context"
	"log"
	"time"
)

func (cfg *apiConfig) StartRefreshTokenCleanup(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			err := cfg.db.DeleteExpiredOrRevokedTokens(context.Background())
			if err != nil {
				log.Printf("ERROR while trying to clean refresh tokens: %v", err)
			}
			log.Printf("Refresh tokens cleaned")
		}
	}()
}
