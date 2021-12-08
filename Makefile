build: ## Build the web app binary
	go build -a -trimpath -o build/api-mocker .

fmt: ## Reformat the code
	go fmt ./...

vet: ## Vet the code
	go vet ./...

test: ## Run the tests
	go test ./... -failfast -coverpkg=./... -coverprofile .testCoverage.txt

coverage: test ## Show test coverage info in the browser
	go tool cover -html .testCoverage.txt

coverage-ci: test ## Shows test coverage in a machine-readable format and creates the cobertura xml file
	command -v gocover-cobertura >> /dev/null || bash -c "pushd . && cd && go install github.com/t-yuki/gocover-cobertura@latest && popd"
	go tool cover -func .testCoverage.txt
	gocover-cobertura < .testCoverage.txt > .testCoverage.xml

