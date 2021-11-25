FROM 10.19.64.203:8080/go/ubuntu_builder AS builder
ARG WORK_DIR

# go build
WORKDIR ${WORK_DIR}
COPY go.* ${WORK_DIR}/
RUN go version && go mod download
COPY . ${WORK_DIR}
RUN CGO_ENABLED=0 go build -ldflags "-w -s -extldflags \"-static\" " -o client
#RUN go build -ldflags "-w -s " -o client

FROM scratch
ARG WORK_DIR

COPY --from=builder ${WORK_DIR}/client /exporter

EXPOSE 9121
ENTRYPOINT [ "/exporter" ]