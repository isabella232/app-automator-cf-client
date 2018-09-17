package internal

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
)

type tokenGetter func() (string, error)

type CapiDoer struct {
    httpClient httpClient
    capiUrl    string
    getToken   tokenGetter
}

func NewCapiDoer(httpClient httpClient, capiUrl string, tokenGetter tokenGetter) *CapiDoer {
    return &CapiDoer{
        httpClient: httpClient,
        capiUrl:    capiUrl,
        getToken:   tokenGetter,
    }
}

type HeaderOption func(header *http.Header)

func (c *CapiDoer) Do(method, path, body string, opts ...HeaderOption) ([]byte, error) {
    req, err := http.NewRequest(method, c.capiUrl+path, ioutil.NopCloser(strings.NewReader(body)))
    if err != nil {
        return nil, err
    }

    req.Header.Add("Content-Type", "application/json")

    for _, o := range opts {
        o(&req.Header)
    }

    if _, ok := req.Header["Authorization"]; !ok {
        token, err := c.getToken()
        if err != nil {
            return nil, err
        }

        req.Header.Add("Authorization", token)
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }

    if code := resp.StatusCode; code > 299 || code < 200 {
        return nil, fmt.Errorf("CAPI request (%s %s) returned unexpected status: %d", method, path, code) //TODO capi error
    }

    if resp.Body == nil {
        return nil, nil
    }
    defer resp.Body.Close()

    return ioutil.ReadAll(resp.Body)
}
