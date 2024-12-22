#!/bin/sh

export ENV_SCHEMA_FILENAME=sample.d/sample.avsc
export ENV_BLOB_KEY=data

jsons2avro() {
	cat sample.d/sample.jsonl |
		ENV_SCHEMA_FILENAME=sample.d/inputgen.avsc \
			json2avrows |
		../avroblobjson2cbor/avroblobjson2cbor |
		cat >./sample.d/input.avro
}

#jsons2avro

export ENV_WASM_MODULES_DIR=./sample.d/modules.d
export ENV_WASM_MODULE_ID=data

cat sample.d/input.avro |
	./avroblob2wasm2blob |
	rq \
		--input-avro \
		--output-json |
	jq \
		-c '.data | .[]' |
	xargs \
		printf '%02x' |
	xxd \
		-plain \
		-revert |
	python3 \
		-m uv \
		tool \
		run \
		cbor2 \
		--sequence
