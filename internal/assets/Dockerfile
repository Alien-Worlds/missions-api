FROM golang

RUN mkdir /home/gocode
RUN mkdir /home/migrations

WORKDIR /home/migrations

RUN go get -u github.com/gobuffalo/packr/v2/...
RUN go get -d -u github.com/gobuffalo/packr/v2/packr2

CMD go mod init "github.com/Alien-Worlds/missions-api" && packr2 clean && packr2 && rm go.mod
