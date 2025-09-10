FROM golang:1.25.1

# Declare build arguments for multi-platform support
ARG TARGETARCH

WORKDIR /usr/


# Download appropriate Tailwind CSS binary based on architecture
RUN case "${TARGETARCH}" in \
        amd64) TAILWIND_ARCH="x64" ;; \
        arm64) TAILWIND_ARCH="arm64" ;; \
        *) echo "Unsupported architecture: ${TARGETARCH}" && exit 1 ;; \
    esac && \
    curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v4.1.11/tailwindcss-linux-${TAILWIND_ARCH} && \
    chmod +x tailwindcss-linux-${TAILWIND_ARCH} && \
    mv tailwindcss-linux-${TAILWIND_ARCH} tailwindcss

RUN curl -sLO https://github.com/saadeghi/daisyui/releases/latest/download/daisyui.js
RUN curl -sLO https://github.com/saadeghi/daisyui/releases/latest/download/daisyui-theme.js

WORKDIR /usr/src/app


RUN go install github.com/air-verse/air@v1.62.0
RUN go install github.com/a-h/templ/cmd/templ@v0.3.943



# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

CMD ["air", "-c", "air.toml"]
