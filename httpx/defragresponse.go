/*
Copyright (c) 2023-2024 Microbus LLC and various contributors

This file and the project encapsulating it are the confidential intellectual property of Microbus LLC.
Neither may be used, copied or distributed without the express written consent of Microbus LLC.
*/

package httpx

import (
	"io"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/microbus-io/fabric/errors"
	"github.com/microbus-io/fabric/frame"
	"github.com/microbus-io/fabric/utils"
)

// DefragResponse merges together multiple fragments back into a single HTTP response
type DefragResponse struct {
	fragments    utils.SyncMap[int, *http.Response]
	maxIndex     atomic.Int32
	count        atomic.Int32
	lastActivity atomic.Int64
}

// NewDefragResponse creates a new response integrator.
func NewDefragResponse() *DefragResponse {
	st := &DefragResponse{}
	st.lastActivity.Store(time.Now().UnixMilli())
	return st
}

// LastActivity indicates how long ago was the last fragment added.
func (st *DefragResponse) LastActivity() time.Duration {
	return time.Duration(time.Now().UnixMilli()-st.lastActivity.Load()) * time.Millisecond
}

// Integrated indicates if all the fragments have been collected and if so returns them as a single HTTP response.
func (st *DefragResponse) Integrated() (integrated *http.Response, err error) {
	maxIndex := st.maxIndex.Load()
	if maxIndex == 1 {
		onlyFrag, _ := st.fragments.Load(1)
		return onlyFrag, nil
	}
	if maxIndex == 0 || st.count.Load() != maxIndex {
		return nil, nil
	}

	// Serialize the bodies of all fragments
	bodies := []io.Reader{}
	var contentLength int64
	for i := 1; i <= int(maxIndex); i++ {
		fragment, ok := st.fragments.Load(i)
		if !ok || fragment == nil {
			return nil, errors.Newf("missing fragment %d", i)
		}
		if fragment.Body == nil {
			return nil, errors.Newf("missing body of fragment %d", i)
		}
		bodies = append(bodies, fragment.Body)
		len, err := strconv.ParseInt(fragment.Header.Get("Content-Length"), 10, 64)
		if err != nil {
			return nil, errors.New("invalid or missing Content-Length header")
		}
		contentLength += len
	}
	integratedBody := io.MultiReader(bodies...)

	// Set the integrated body on the first fragment
	firstFragment, ok := st.fragments.Load(1)
	if !ok || firstFragment == nil {
		return nil, errors.New("missing first fragment")
	}
	frame.Of(firstFragment).SetFragment(1, 1) // Clear the header
	firstFragment.Header.Set("Content-Length", strconv.FormatInt(contentLength, 10))
	firstFragment.Body = io.NopCloser(integratedBody)
	return firstFragment, nil
}

// Add a fragment to be integrated.
func (st *DefragResponse) Add(r *http.Response) error {
	index, max := frame.Of(r).Fragment()
	st.maxIndex.Store(int32(max))
	st.fragments.Store(index, r)
	st.count.Add(1)
	st.lastActivity.Store(time.Now().UnixMilli())
	return nil
}
