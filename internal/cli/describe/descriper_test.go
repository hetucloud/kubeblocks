/*
Copyright ApeCloud Inc.

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

package describe

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	cmdtesting "k8s.io/kubectl/pkg/cmd/testing"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/describe"

	"github.com/apecloud/kubeblocks/internal/cli/cluster"
	"github.com/apecloud/kubeblocks/internal/cli/testing"
	"github.com/apecloud/kubeblocks/internal/cli/types"
)

var _ = Describe("Describer", func() {
	tf := cmdtesting.NewTestFactory().WithNamespace("test")
	defer tf.Cleanup()

	It("describer map", func() {
		describer, err := DescriberFn(tf, &meta.RESTMapping{
			Resource:         types.ClusterGVR(),
			GroupVersionKind: types.ClusterGK().WithVersion(types.Version),
		})
		Expect(describer).ShouldNot(BeNil())
		Expect(err).ShouldNot(HaveOccurred())

		describer, err = DescriberFn(tf, &meta.RESTMapping{
			Resource: schema.GroupVersionResource{
				Group:    types.Group,
				Version:  types.Version,
				Resource: "tests",
			},
			GroupVersionKind: schema.GroupVersionKind{
				Group:   types.Group,
				Version: types.Version,
				Kind:    "test",
			}})
		Expect(describer).ShouldNot(BeNil())
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("describe cluster", func() {
		It("describe return error", func() {
			describer := &ClusterDescriber{
				client: testing.FakeClientSet(&corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "foo",
					},
					Spec: corev1.PodSpec{
						ServiceAccountName: "fooaccount",
					},
				}),
				dynamic: testing.FakeDynamicClient(),
			}
			describerSettings := describe.DescriberSettings{ShowEvents: true, ChunkSize: cmdutil.DefaultChunkSize}
			res, err := describer.Describe("test", "test", describerSettings)
			Expect(res).Should(Equal(""))
			Expect(err).Should(HaveOccurred())
		})

		It("mock cluster and check", func() {
			describer := ClusterDescriber{ClusterObjects: cluster.FakeClusterObjs()}
			res, err := describer.describeCluster(nil)
			Expect(res).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
