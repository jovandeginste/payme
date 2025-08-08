release-patch release-minor release-major:
	$(MAKE) release VERSION=$(shell go run github.com/mdomke/git-semver/v6@latest -target $(subst release-,,$@))

release:
	git tag -s -a $(VERSION) -m "Release $(VERSION)"
	@echo "Now run:"
	@echo "- git push --tags"
	@echo "- gh release create --generate-notes $(VERSION)"
