#Prebuild node
FROM alpine/git AS prebuild

WORKDIR /

RUN git clone https://github.com/Dilshat/telegram-bot.git 


#Builder node
FROM golang:1.13.4-alpine3.10 as builder

WORKDIR /

RUN mkdir telegram-bot

COPY --from=prebuild /telegram-bot /telegram-bot 

WORKDIR telegram-bot

RUN GOINSECURE=1 go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o tbot

#Production node
FROM alpine:latest AS production

WORKDIR /

#copy .env containing default values from builder node since specific env vars are to be supplied on docker container startup
COPY --from=builder /telegram-bot/.env .
COPY --from=builder /telegram-bot/tbot .

#copy scripts from host
RUN mkdir scripts
COPY scripts/* scripts/

#copy attachments from host
RUN mkdir attachments
COPY attachments/* attachments/

CMD ["./tbot"]


