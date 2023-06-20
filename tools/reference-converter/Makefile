.PHONY: build

build:
	docker buildx build -t reference-converter:builder --target builder .
	docker buildx build -t reference-converter .

/results:
	@mkdir -p $(PWD)/results

test: /results
	docker run --rm -v "$(PWD)/results:/results" \
		reference-converter:builder ./test.sh /results/test.xml
