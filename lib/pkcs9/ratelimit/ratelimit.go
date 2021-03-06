// Copyright © SAS Institute Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ratelimit

import (
	"context"
	"log"
	"time"

	"github.com/sassoftware/relic/lib/pkcs7"
	"github.com/sassoftware/relic/lib/pkcs9"
	"golang.org/x/time/rate"
)

type limiter struct {
	Timestamper pkcs9.Timestamper
	Limit       *rate.Limiter
}

func New(t pkcs9.Timestamper, r float64, burst int) pkcs9.Timestamper {
	if r == 0 {
		return t
	}
	if burst < 1 {
		burst = 1
	}
	return &limiter{t, rate.NewLimiter(rate.Limit(r), burst)}
}

func (l *limiter) Timestamp(ctx context.Context, req *pkcs9.Request) (*pkcs7.ContentInfoSignedData, error) {
	start := time.Now()
	if err := l.Limit.Wait(ctx); err != nil {
		return nil, err
	}
	if waited := time.Now().Sub(start); waited > 50*time.Millisecond {
		log.Printf("timestamper: waited %s due to rate limit", waited)
	}
	return l.Timestamper.Timestamp(ctx, req)
}
