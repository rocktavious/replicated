package platformclient_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/replicatedhq/replicated/pkg/platformclient"
	"net/http"
)

var _ = Describe("Channel", func() {
	Context("CreateChannel", func() {
		var (
			server        *ghttp.Server
			appID         = "some-app-id"
			apiKey        = "some-api-key"
			basicRespBody = `
[{
	"Id": "some-id",
    "Name": "some-channel-name"
}]
`
		)

		BeforeEach(func() {
			server = ghttp.NewServer()
		})

		AfterEach(func() {
			server.Close()
		})

		It("creates a channel successfully", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", fmt.Sprintf("/v1/app/%s/channel", appID)),
					ghttp.VerifyHeader(http.Header{
						"Authorization": []string{apiKey},
					}),

					ghttp.RespondWith(http.StatusOK, basicRespBody),
				),
			)

			client := platformclient.NewHTTPClient(server.URL(), apiKey)
			err := client.CreateChannel("some-app-id", "some-name", "some-description")
			Expect(err).ToNot(HaveOccurred())
		})

		It("propagates errors from DoJSON ", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", fmt.Sprintf("/v1/app/%s/channel", appID)),
					ghttp.VerifyHeader(http.Header{
						"Authorization": []string{apiKey},
					}),

					ghttp.RespondWith(http.StatusNotFound, basicRespBody),
				),
			)

			client := platformclient.NewHTTPClient(server.URL(), apiKey)
			err := client.CreateChannel("some-app-id", "some-name", "some-description")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(platformclient.ErrNotFound.Error()))
		})
	})
})
