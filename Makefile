help:
	echo "asdsad"

format:
	@gofmt -s -w .

review: format
	@echo "============= Spell Check ============= "
	@misspell .
	
	@echo "============= Ineffectual Assignments Check ============= "
	@ineffassign .

	@echo "============= Duplication Check ============= "
	find . -not -name '*_test.go' -name '*.go' | dupl -t 30 -files

	@echo "============= Repeated Strings Check ============= "
	@goconst .

	@echo "============= Security Check ============= "
	@gosec .

	@echo "============= Vet Check ============= "
	@go vet --all .

	@echo "============= Preallocation Check ============= "
	@prealloc -forloops -set_exit_status -simple -rangeloops .

	@echo "============= Shadow Variables Check ============= "
	@shadow -strict .

	@echo "============= Cyclomatic Complexity Check ============= "
	@gocyclo -total -ignore "_test" -over 8 -avg .

tools:
	go install github.com/client9/misspell/cmd/misspell@latest
	go install github.com/gordonklaus/ineffassign@latest
	go install github.com/mibk/dupl@latest
	go install github.com/jgautheron/goconst/cmd/goconst@latest
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install github.com/alexkohler/prealloc@latest
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest