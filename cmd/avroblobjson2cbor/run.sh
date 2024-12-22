#!/bin/sh

export ENV_SCHEMA_FILENAME=sample.d/sample.avsc

jsons2avro(){
	cat sample.d/sample.jsonl |
    	ENV_SCHEMA_FILENAME=sample.d/inputgen.avsc \
			json2avrows |
		cat > ./sample.d/input.avro
}

#jsons2avro

export ENV_BLOB_KEY=data
export ENV_CONVERT_BYTES_TO_STRING=true

cat sample.d/input.avro |
	./avroblobjson2cbor |
	rq -a -J |
	jq -c '.data|.[]' |
	xargs printf '%02x' |
	xxd \
		-revert \
		-plain |
	python3 \
		-m uv \
		tool \
		run \
		cbor2 \
			--sequence
