.PHONY: bump

bump:
	@VERSION=$$(go run . up -w -f main.go -p 'version\s*=\s*"(\d+\.\d+\.\d+)"') && \
		[ -n "$$VERSION" ] && \
		git commit -am "bump version to $$VERSION" && \
		git tag "v$$VERSION" && \
		git push origin "v$$VERSION"
