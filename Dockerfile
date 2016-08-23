FROM scratch
ADD quicklog-linux /quicklog
ADD ql2etcd-linux /ql2etcd
ADD qlsearch-linux /qlsearch
ADD static /static
WORKDIR /

CMD ["/quicklog"]

