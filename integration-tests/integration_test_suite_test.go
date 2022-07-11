package integration_test

import (
	"encoding/json"
	"testing"

	testframework "github.com/cloudfoundry/cloud-service-broker/brokerpaktestframework"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIntegrationTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IntegrationTests Suite")
}

const (
	awsSecretAccessKey = "aws-secret-access-key"
	awsAccessKeyID     = "aws-access-key-id"
)

var (
	mockTerraform testframework.TerraformMock
	broker        *testframework.TestInstance
)
// TODO credentials are hardcoded
// TODO invocation vars include the command
// TODO used the broker's terraform and not the system one
var _ = BeforeSuite(func() {
	var err error
	mockTerraform, err = testframework.NewTerraformMock()
	Expect(err).NotTo(HaveOccurred())

	broker, err = testframework.BuildTestInstance(testframework.PathToBrokerPack(), mockTerraform, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	Expect(broker.Start(GinkgoWriter, []string{
		"GSB_SERVICE_CSB_AWS_S3_BUCKET_PLANS=" + marshall(customS3Plans),
		"AWS_ACCESS_KEY_ID=",
		"AWS_SECRET_ACCESS_KEY=",
		"CSB_LISTENER_HOST=localhost",
		"GSB_COMPATIBILITY_ENABLE_BETA_SERVICES=true",
	})).To(Succeed())
})

var _ = AfterSuite(func() {
	if broker != nil {
		Expect(broker.Cleanup()).To(Succeed())
	}
})

func marshall(element any) string {
	b, err := json.Marshal(element)
	Expect(err).NotTo(HaveOccurred())
	return string(b)
}
