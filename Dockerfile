FROM ubuntu

COPY licenseserver /licenseserver
COPY conf /conf
COPY views /views
COPY public /public

EXPOSE 8080

WORKDIR /
CMD ["/licenseserver"]
