FROM node:latest

RUN mkdir -p app/

WORKDIR /app

COPY wicigga_frontend/package*.json ./

RUN npm install

COPY wicigga_frontend/ ./ 

RUN npm run build