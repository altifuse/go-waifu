FROM golang:1.22-bookworm

# install waifu2x
RUN apt update && \
    DEBIAN_FRONTEND=noninteractive apt install libvulkan-dev libgomp1 unzip -y
ADD https://github.com/nihui/waifu2x-ncnn-vulkan/releases/download/20220728/waifu2x-ncnn-vulkan-20220728-ubuntu.zip /tmp/waifu2x.zip
RUN unzip /tmp/waifu2x.zip -d /
RUN mv /waifu2x-ncnn-vulkan-20220728-ubuntu /waifu2x
RUN rm /tmp/waifu2x.zip

# compile go-waifu
ADD api /tmp/go-waifu
WORKDIR /tmp/go-waifu
RUN go build
RUN mv go-waifu /usr/bin/go-waifu
WORKDIR /
RUN rm -r /tmp/go-waifu

ENTRYPOINT [ "go-waifu" ]