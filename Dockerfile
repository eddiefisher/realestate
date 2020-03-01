FROM scratch

RUN mkdir /config

COPY build/realestate /
COPY config/* /config

CMD ["/realestate"]