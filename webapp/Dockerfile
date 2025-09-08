FROM node:24-alpine AS builder

WORKDIR /app

COPY package.json pnpm-lock.yaml ./

RUN npm install -g pnpm

RUN pnpm install

COPY . .

RUN pnpm run build

FROM node:24-alpine AS runner

WORKDIR /app

COPY --from=builder /app/.output .output

EXPOSE 3000

CMD ["pnpm", "run", "preview"]
