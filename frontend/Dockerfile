FROM node:14
USER root
RUN mkdir /app
COPY . /app
WORKDIR /app/
# RUN apt-get add --update nodejs npm
RUN npm i -g parcel-bundler
RUN npm install
ENTRYPOINT ["parcel", "--hmr-port", "7001", "index.html"]
EXPOSE 1234
