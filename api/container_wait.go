package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/docker/docker/api/types/container"
)

// ContainerWait waits until the specified container is in a certain state
// indicated by the given condition, either "not-running" (default),
// "next-exit", or "removed".
//
// If this client's API version is before 1.30, condition is ignored and
// ContainerWait will return immediately with the two channels, as the server
// will wait as if the condition were "not-running".
//
// If this client's API version is at least 1.30, ContainerWait blocks until
// the request has been acknowledged by the server (with a response header),
// then returns two channels on which the caller can wait for the exit status
// of the container or an error if there was a problem either beginning the
// wait request or in getting the response. This allows the caller to
// synchronize ContainerWait with other calls, such as specifying a
// "next-exit" condition before issuing a ContainerStart request.
func (api *API) ContainerWait(
	ctx context.Context,
	containerID string,
	condition container.WaitCondition,
	timeout time.Duration,
) (<-chan container.ContainerWaitOKBody, <-chan error) {
	resultC := make(chan container.ContainerWaitOKBody)
	errC := make(chan error, 1)

	query := url.Values{}
	query.Set("condition", string(condition))

	path := fmt.Sprintf("/containers/%s/wait?%s", containerID, query.Encode())
	url := fmt.Sprintf("%s%s", api.endpoint, path)
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		errC <- err
		return resultC, errC
	}

	client := &http.Client{Timeout: timeout}
	if _, err = client.Do(request); err != nil {
		errC <- err
		return resultC, errC
	}

	go func() {
		resp, err := client.Do(request)
		if err != nil {
			errC <- err
		}

		defer resp.Body.Close()

		var res container.ContainerWaitOKBody
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errC <- err
		}

		if err := json.Unmarshal(body, &res); err != nil {
			errC <- err
		}

		resultC <- res
	}()
	return resultC, errC
}
