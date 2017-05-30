FROM ubuntu

COPY licenseserver /licenseserver
COPY conf /conf

EXPOSE 8080

WORKDIR /
CMD ["/licenseserver"]
