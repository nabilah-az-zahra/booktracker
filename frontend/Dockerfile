FROM node:24-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

ARG VITE_GOOGLE_BOOKS_API_KEY
ENV VITE_GOOGLE_BOOKS_API_KEY=$VITE_GOOGLE_BOOKS_API_KEY

RUN npm run build

FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

RUN touch /var/run/nginx.pid && \
    chown -R nginx:nginx /var/run/nginx.pid /var/cache/nginx /var/log/nginx /usr/share/nginx/html
USER nginx

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]