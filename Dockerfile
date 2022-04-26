FROM scratch

ADD rootfs/ /

CMD ["/bin/ash"]