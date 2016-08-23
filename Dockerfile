FROM scratch
ADD bin/quicklog-linux /quicklog
WORKDIR /

CMD ["/quicklog"]

