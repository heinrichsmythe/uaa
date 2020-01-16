package k8s_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"path/filepath"
)

const UaaYmlConfigKey = "uaa.yml"

var _ = Describe("Uaa ConfigMap", func() {
	var configPath, uaaLibPath, valuesPath string
	var database Database

	BeforeEach(func() {
		configPath = pathToTemplate("config.yml")
		uaaLibPath = pathToTemplate("uaa.lib.yml")
		valuesPath = pathToTemplate(filepath.Join("values", "values.yml"))
		database = Database{Username: "uaa", Password: "password", Url: "http://example.com"}
	})

	It("Renders a config map with default values", func() {
		ctx := NewRenderingContext(configPath, uaaLibPath, valuesPath)

		Expect(ctx).To(
			ProduceYAML(
				RepresentingConfigMap().WithDataFieldMatching(UaaYmlConfigKey, func(uaaYml *DataFieldMatcher) {
					uaaYml.WithFields(Fields{
						"LoginSecret": Equal("loginsecret"),
						"Issuer":      Equal(Issuer{Uri: "http://localhost:8080/uaa"}),
						"Database": MatchFields(IgnoreExtras, Fields{
							"Username": Equal(database.Username),
							"Password": Equal(database.Password),
							"Url":      Equal(database.Url),
						}),
					})
				}),
			))
	})

	It("Can renders a config map with overriden values", func() {
		database.Username = "database-username"
		database.Password = "database-password"

		ctx := NewRenderingContext(configPath, uaaLibPath, valuesPath).WithData(map[string]string{
			"database.username": database.Username,
			"database.password": database.Password,
		})

		Expect(ctx).To(
			ProduceYAML(
				RepresentingConfigMap().WithDataFieldMatching(UaaYmlConfigKey, func(uaaYml *DataFieldMatcher) {
					uaaYml.WithFields(Fields{
						"LoginSecret": Equal("loginsecret"),
						"Issuer":      Equal(Issuer{Uri: "http://localhost:8080/uaa"}),
						"Database": MatchFields(IgnoreExtras, Fields{
							"Username": Equal(database.Username),
							"Password": Equal(database.Password),
							"Url":      Equal(database.Url),
						}),
					})
				}),
			))
	})
})
