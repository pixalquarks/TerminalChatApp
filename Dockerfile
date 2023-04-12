FROM ubuntu
COPY /dist .
COPY .env .
EXPOSE 80
CMD ./dist