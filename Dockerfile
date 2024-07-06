FROM gcr.io/distroless/cc

ADD app /app

EXPOSE 8080
CMD ["/app"]
