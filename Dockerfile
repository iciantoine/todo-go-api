FROM golang:latest as build

COPY . /
WORKDIR /
RUN make dist-docker prepare

# Prepare final, minimal image
FROM heroku/heroku:22

COPY --from=build /dist /app
COPY --from=build /go/bin/tern /app/

ENV HOME /app
WORKDIR /app
RUN useradd -m heroku
USER heroku
CMD /app/linux/server
