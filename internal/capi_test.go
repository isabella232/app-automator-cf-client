package internal_test

import (
    . "github.com/onsi/ginkgo"
    "github.com/pivotal-cf/eats-cf-client/internal"
    . "github.com/onsi/gomega"
    . "github.com/onsi/ginkgo/extensions/table"
    "net/http"
    "github.com/pkg/errors"
    "github.com/pivotal-cf/eats-cf-client/models"
)

var _ = Describe("Capi", func() {
    Describe("Apps()", func() {
        It("gets the apps", func() {
            c := internal.NewCapiClient(func(method, path string, body string) ([]byte, error) {
                Expect(method).To(Equal(http.MethodGet))
                Expect(path).To(And(
                    ContainSubstring("/v3/apps"),
                    ContainSubstring("lemons=limes"),
                    ContainSubstring("mangoes=limes"),
                ))
                return []byte(validAppsResponse), nil
            })

            apps, err := c.Apps(map[string]string{
                "lemons":  "limes",
                "mangoes": "limes",
            })

            Expect(err).ToNot(HaveOccurred())
            Expect(apps).To(ConsistOf(models.App{
                Guid: "app-guid",
            }))
        })

        DescribeTable("errors", func(do func(method, path string, body string) ([]byte, error)) {
            c := internal.NewCapiClient(do)

            _, err := c.Apps(nil)
            Expect(err).To(HaveOccurred())
        },
            Entry("do returns an error", func(method, path string, body string) ([]byte, error) {
                return nil, errors.New("expected")
            }),
            Entry("do returns invalid json", func(method, path string, body string) ([]byte, error) {
                return []byte("{]"), nil
            }),
        )
    })

    Describe("Process()", func() {
        It("gets the process", func() {
            c := internal.NewCapiClient(func(method, path string, body string) ([]byte, error) {
                Expect(method).To(Equal(http.MethodGet))
                Expect(path).To(Equal("/v3/apps/app-guid/processes/process-type"))
                return []byte(validProcessResponse), nil
            })

            process, err := c.Process("app-guid", "process-type")
            Expect(err).ToNot(HaveOccurred())

            Expect(process).To(Equal(models.Process{
                Instances: 2,
            }))
        })

        DescribeTable("errors", func(do func(method, path string, body string) ([]byte, error)) {
            c := internal.NewCapiClient(do)

            _, err := c.Process("app-guid", "process-type")
            Expect(err).To(HaveOccurred())
        },
            Entry("do returns an error", func(method, path string, body string) ([]byte, error) {
                return nil, errors.New("expected")
            }),
            Entry("do returns invalid json", func(method, path string, body string) ([]byte, error) {
                return []byte("{]"), nil
            }),
        )
    })

    Describe("Scale()", func() {
        It("scales the process", func() {
            var called bool
            c := internal.NewCapiClient(func(method, path string, body string) ([]byte, error) {
                called = true
                Expect(method).To(Equal(http.MethodPost))
                Expect(path).To(Equal("/v3/apps/app-guid/processes/process-type/actions/scale"))
                Expect(body).To(MatchJSON(`{ "instances": 5 }`))
                return nil, nil
            })

            err := c.Scale("app-guid", "process-type", 5)
            Expect(err).ToNot(HaveOccurred())
            Expect(called).To(BeTrue())
        })

        DescribeTable("errors", func(do func(method, path string, body string) ([]byte, error)) {
            c := internal.NewCapiClient(do)
            Expect(c.Scale("app-guid", "process-type", 5)).ToNot(Succeed())
        },
            Entry("do returns an error", func(method, path string, body string) ([]byte, error) {
                return nil, errors.New("expected")
            }),
        )
    })
})

const validAppsResponse = `{"resources": [ { "guid": "app-guid" } ]}`
const validProcessResponse = `{ "instances": 2 }`