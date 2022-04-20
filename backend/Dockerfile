# DEVELOPMENT ENVIRONMENT
FROM node:16.14-slim as dev
WORKDIR /opt/jacobianmatthews.com/backend
COPY . .
RUN yarn install
ENV NODE_ENV=development
CMD ["yarn", "start:dev"]

# TESTING ENVIRONMENT
FROM node:16.14-slim as test
WORKDIR /opt/jacobianmatthews.com/backend
COPY . .
RUN yarn install --frozen-lockfile --non-interactive
RUN yarn audit --groups dependencies
RUN yarn lint
RUN yarn test
RUN yarn test:e2e

# BUILD ENVIRONMENT
FROM node:16-14-slim as build
WORKDIR /opt/jacobianmatthews.com/backend
COPY . .
RUN yarn install --frozen-lockfile --non-interactive
RUN yarn build

# PRODUCTION ENVIRONMENT
FROM node:16-14-slim as prod
WORKDIR /opt/jacobianmatthews.com/backend
COPY --from=build /opt/jacobianmatthews.com/backend/package.json .
COPY --from=build /opt/jacobianmatthews.com/backend/yarn.lock .
COPY --from=build /opt/jacobianmatthews.com/backend/dist .
ENV NODE_ENV=production
RUN yarn install --frozen-lockfile --non-interactive --production
CMD ["yarn", "start:prod"]