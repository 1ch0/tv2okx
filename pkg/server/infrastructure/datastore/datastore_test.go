package datastore

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/1ch0/tv2okx/pkg/server/domain/model"
)

var _ = Describe("Test new entity function", func() {

	It("Test new application entity", func() {
		var app model.SystemInfo
		new, err := NewEntity(&app)
		Expect(err).To(BeNil())
		err = json.Unmarshal([]byte(`{"name":"demo"}`), new)
		Expect(err).To(BeNil())
		diff := cmp.Diff(new.PrimaryKey(), "demo")
		Expect(diff).Should(BeEmpty())
	})

	It("Test new multiple application entity", func() {
		var app model.SystemInfo
		var list []Entity
		var n = 3
		for n > 0 {
			new, err := NewEntity(&app)
			Expect(err).To(BeNil())
			err = json.Unmarshal([]byte(fmt.Sprintf(`{"name":"demo %d"}`, n)), new)
			Expect(err).To(BeNil())
			diff := cmp.Diff(new.PrimaryKey(), fmt.Sprintf("demo %d", n))
			Expect(diff).Should(BeEmpty())
			list = append(list, new)
			n--
		}
		diff := cmp.Diff(list[0].PrimaryKey(), "demo 3")
		Expect(diff).Should(BeEmpty())
	})

})
