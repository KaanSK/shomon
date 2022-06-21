FROM scratch
COPY shomon /usr/local/bin/shomon
ENTRYPOINT [ "/usr/local/bin/shomon" ]