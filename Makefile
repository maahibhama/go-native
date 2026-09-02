.PHONY: test vet benchmark ios run-ios android run-android
test:
	GOCACHE=/tmp/go-native-gocache go test ./...
vet:
	GOCACHE=/tmp/go-native-gocache go vet ./...
benchmark:
	GOCACHE=/tmp/go-native-gocache go test -bench=. -benchmem ./runtime
ios:
	GOCACHE=/tmp/go-native-gocache ./scripts/build-ios.sh
run-ios:
	GOCACHE=/tmp/go-native-gocache ./scripts/run-ios.sh
android:
	GOCACHE=/tmp/go-native-gocache GOPATH=/tmp/go-native-gopath ./scripts/build-android.sh
run-android:
	GOCACHE=/tmp/go-native-gocache GOPATH=/tmp/go-native-gopath ./scripts/run-android.sh
