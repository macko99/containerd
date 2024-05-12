/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package containerd

import (
	"context"
	"time"

	"github.com/containerd/containerd/leases"
	"k8s.io/klog/v2"
)

// WithLease attaches a lease on the context
func (c *Client) WithLease(ctx context.Context, opts ...leases.Opt) (context.Context, func(context.Context) error, error) {
	klog.Infof("%s [CONTINUUM] 0970 containerd:WithLease:start sandbox=%s", time.Now().UnixNano(), ctx)
	nop := func(context.Context) error { return nil }

	_, ok := leases.FromContext(ctx)
	if ok {
		return ctx, nop, nil
	}

	ls := c.LeasesService()

	if len(opts) == 0 {
		// Use default lease configuration if no options provided
		opts = []leases.Opt{
			leases.WithRandomID(),
			leases.WithExpiration(24 * time.Hour),
		}
	}

	l, err := ls.Create(ctx, opts...)
	if err != nil {
		return ctx, nop, err
	}

	ctx = leases.WithLease(ctx, l.ID)
	klog.Infof("%s [CONTINUUM] 0971 containerd:WithLease:done sandbox=%s", time.Now().UnixNano(), ctx)
	return ctx, func(ctx context.Context) error {
		return ls.Delete(ctx, l)
	}, nil
}
