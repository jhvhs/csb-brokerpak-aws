package integration_test

import (
	"encoding/json"
	"fmt"
	testframework "github.com/cloudfoundry/cloud-service-broker/brokerpaktestframework"
	tfjson "github.com/hashicorp/terraform-json"

	//"fmt"
	//
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("S3", Label("s3"), func() {
	const s3ServiceName = "csb-aws-s3-bucket"
	BeforeEach(func() {
		Expect(mockTerraform.SetTFState([]testframework.TFStateValue{})).To(Succeed())
	})

	AfterEach(func() {
		Expect(mockTerraform.Reset()).To(Succeed())
	})

	FIt("should parse the hcl", func() {
		_, err := broker.Provision(s3ServiceName, "private", nil)

		//Expect(err).NotTo(HaveOccurred())
		//Expect(mockTerraform.FirstTerraformInvocationVars()).To(
		//	SatisfyAll(
		//		HaveKeyWithValue("bucket_name", "csb-"+instanceID),
		//		HaveKeyWithValue("enable_versioning", false),
		//		HaveKeyWithValue("labels", HaveKeyWithValue("pcf-instance-id", instanceID)),
		//		HaveKeyWithValue("region", "us-west-2"),
		//		HaveKeyWithValue("acl", "private"),
		//		HaveKeyWithValue("aws_access_key_id", awsAccessKeyID),
		//		HaveKeyWithValue("aws_secret_access_key", awsSecretAccessKey),
		//	),
		//)



		var outputJson tfjson.Plan


		out, err := mockTerraform.FirstTerraformInvocationPlan()
		Expect(err).ToNot(HaveOccurred())
		err = json.Unmarshal([]byte(out), &outputJson)
		Expect(err).NotTo(HaveOccurred())

		for _, change := range outputJson.ResourceChanges {
			fmt.Printf("Change Resource Name: %s Type:%s Change: %s", change.Name, change.Type, change.Change)
			Expect(change.Type).To(Equal("aws_s3_bucket"))
			Expect(change.Name).To(Equal("b"))
			Expect(change.Change.Actions).To(ConsistOf(tfjson.Action("create")))
			Expect(change.Change.After).To(MatchAllKeys(Keys{"bucket": Equal("somename"), "bucket_prefix": BeNil(), "force_destroy": BeFalse(), "tags": Not(BeEmpty()), "tags_all": Not(BeEmpty())}))
		}

	})

})

// TODO
// how do we tie this to the broker? can we hook into what the broker executes with another terraform mock that does this under an apply?
// pros: we would get exactly what the broker would execute, including variables and workspace. We wouldn't need to do anything to the path
// cons: more convoluted
// Do we need to tie to the broker or should we just run tf?
// cons: this looks more shady
// pros: simples
//
//
//func ShowPlan(args ...string) *gexec.Session {
//	fmt.Fprintf(GinkgoWriter, "Running: tf plan %s\n", strings.Join(args, " "))
//
//	command := exec.Command("terraform", "-chdir=../terraform/s3/provision", "init")
//	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
//	Expect(err).NotTo(HaveOccurred())
//	Eventually(session, 5*time.Minute).Should(gexec.Exit(0))
//	tmpFile, _ := ioutil.TempFile("", "test-tf")
//	allArgs := []string{"-chdir=../terraform/s3/provision", "plan", "-refresh=false",  fmt.Sprintf("-out=%s", tmpFile.Name())}
//	allArgs = append(allArgs, args...)
//	command = exec.Command("terraform", allArgs...)
//	session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
//	Expect(err).NotTo(HaveOccurred())
//	Eventually(session, 5*time.Minute).Should(gexec.Exit(0))
//	command = exec.Command("terraform", "-chdir=../terraform/s3/provision", "show", "-json", tmpFile.Name())
//	session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
//	Expect(err).NotTo(HaveOccurred())
//	Eventually(session, 5*time.Minute).Should(gexec.Exit(0))
//	return session
//}

