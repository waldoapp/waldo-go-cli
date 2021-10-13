# This creates a minimal Docker image to allow easy testing of Waldo Go CLI on
# Linux. It is based on the Alpine Linux image. It also installs the `waldo`
# executable into `/usr/local/bin`.
#
# To build the image, issue the following commands from the project directory:
#
# ```
# $ make build_linux
# $ docker build -t linux-waldo-go-cli .
# ```
#
# You can subsequently run the resulting Docker image from any directory (and
# also create a readonly binding to that directory) with the following command:
#
# ```
# $ docker run -it -v $(pwd):/app:ro linux-waldo-go-cli
# ```

FROM alpine:latest

COPY bin/waldo-linux-amd64 /usr/local/bin/waldo

RUN chmod a+x /usr/local/bin/waldo
