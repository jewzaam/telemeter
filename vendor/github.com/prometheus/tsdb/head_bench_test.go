// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tsdb

import (
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/prometheus/tsdb/labels"
	"github.com/prometheus/tsdb/testutil"
)

func BenchmarkHeadStripeSeriesCreate(b *testing.B) {
	// Put a series, select it. GC it and then access it.
	h, err := NewHead(nil, nil, nil, 1000)
	testutil.Ok(b, err)
	defer h.Close()

	for i := 0; i < b.N; i++ {
		h.getOrCreate(uint64(i), labels.FromStrings("a", strconv.Itoa(i)))
	}
}

func BenchmarkHeadStripeSeriesCreateParallel(b *testing.B) {
	// Put a series, select it. GC it and then access it.
	h, err := NewHead(nil, nil, nil, 1000)
	testutil.Ok(b, err)
	defer h.Close()

	var count int64

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := atomic.AddInt64(&count, 1)
			h.getOrCreate(uint64(i), labels.FromStrings("a", strconv.Itoa(int(i))))
		}
	})
}

// TODO: generalize benchmark and pass all postings for matchers here
func BenchmarkHeadPostingForMatchers(b *testing.B) {
	// Put a series, select it. GC it and then access it.
	h, err := NewHead(nil, nil, nil, 1000)
	testutil.Ok(b, err)
	defer h.Close()

	// TODO: vary number of series
	for i := 0; i < 100; i++ {
		h.getOrCreate(uint64(i), labels.FromStrings("a", strconv.Itoa(i)))
	}

	b.ResetTimer()

	all, _ := labels.NewRegexpMatcher("a", ".*")

	for i := 0; i < b.N; i++ {
		_, err := PostingsForMatchers(h.indexRange(0, 1000), all)
		testutil.Ok(b, err)
	}
}
