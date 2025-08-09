FROM node:24-slim AS base
WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Install dependencies
COPY package.json ./
COPY pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# Prepare application
COPY . ./
RUN pnpm build

CMD ["pnpm", "preview", "--host", "'0.0.0.0'", "--port", "5173"]