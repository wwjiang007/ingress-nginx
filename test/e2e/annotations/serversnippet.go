/*
Copyright 2018 The Kubernetes Authors.

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

package annotations

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/ingress-nginx/test/e2e/framework"
	"strings"
)

var _ = framework.IngressNginxDescribe("Annotations - ServerSnippet", func() {
	f := framework.NewDefaultFramework("serversnippet")

	BeforeEach(func() {
		err := f.NewEchoDeploymentWithReplicas(2)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
	})

	It(`add valid directives to server via server snippet"`, func() {
		host := "serversnippet.foo.com"
		annotations := map[string]string{
			"nginx.ingress.kubernetes.io/server-snippet": `
				more_set_headers "Content-Length: $content_length";
				more_set_headers "Content-Type: $content_type";`,
		}

		ing := framework.NewSingleIngress(host, "/", host, f.IngressController.Namespace, "http-svc", 80, &annotations)
		_, err := f.EnsureIngress(ing)

		Expect(err).NotTo(HaveOccurred())
		Expect(ing).NotTo(BeNil())

		err = f.WaitForNginxServer(host,
			func(server string) bool {
				return strings.Contains(server, `more_set_headers "Content-Length: $content_length`) && strings.Contains(server, `more_set_headers "Content-Type: $content_type";`)
			})
		Expect(err).NotTo(HaveOccurred())
	})
})
