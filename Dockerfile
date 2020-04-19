FROM golang:1.13
EXPOSE 80

ARG DEPLOY_ENV=production

ADD src /opt/joonas.ninja-chat/
RUN cp -R ./env/${DEPLOY_ENV}.env /opt/joonas.ninja-chat/app.env

WORKDIR /opt/joonas.ninja-chat
RUN go build

CMD go run *.go


