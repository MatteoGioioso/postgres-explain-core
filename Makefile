.SILENT:

install-tools:
	@cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

generate-types:
	tygo generate

parse:
	ts-node components/parser.ts | jq '.'