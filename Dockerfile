FROM scratch
COPY main /
ENTRYPOINT ["./main"]
EXPOSE 8080