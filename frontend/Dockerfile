FROM 1.117.192.82:8666/library/node:20 AS builder

WORKDIR /workspace
COPY package.json yarn.lock ./
RUN --mount=type=cache,target=/usr/local/share/.cache,id=yarn-hgt-cache \
    yarn install

COPY . .
RUN yarn build && rm -rf node_modules/.cache


FROM 1.117.192.82:8666/library/nginx:1.27-alpine
COPY --from=builder /workspace/dist /dist

EXPOSE 8840
EXPOSE 443