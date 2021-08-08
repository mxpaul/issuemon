FROM golang:1.16.7
WORKDIR /src/issuemon

build:
	COPY . .
	RUN ./build/build_binaries.sh
	SAVE ARTIFACT dist_build/bin AS LOCAL dist_build/bin

test:
	COPY . .
	ARG GOFLAGS
	RUN make GOFLAGS="${GOFLAGS}" test
