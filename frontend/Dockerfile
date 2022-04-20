# DEVELOPMENT ENVIRONMENT
FROM node:16.14-slim as dev
WORKDIR /opt/jacobianmatthews.com/frontend
COPY . .
RUN yarn install
ENV NODE_ENV=development
CMD ["yarn", "dev"]

# TEST ENVIRONMENT
FROM node:16.14-slim as test
WORKDIR /opt/jacobianmatthews.com/frontend
COPY . .
RUN yarn install --frozen-lockfile --non-interactive
RUN yarn audit --groups dependencies
RUN yarn lint
RUN yarn test:ci

# BUILD ENVIRONMENT
FROM node:16.14-slim as build
WORKDIR /opt/jacobianmatthews.com/frontend
COPY . .
RUN yarn install --frozen-lockfile --non-interactive
RUN yarn build

# PRODUCTION ENVIRONMENT 
FROM node:16.14-slim as prod
WORKDIR /opt/jacobianmatthews.com/frontend

COPY --from=build /opt/jacobianmatthews.com/frontend/package.json .
COPY --from=build /opt/jacobianmatthews.com/frontend/yarn.lock .
COPY --from=build /opt/jacobianmatthews.com/frontend/.next .
COPY --from=build /opt/jacobianmatthews.com/frontend/public .

ENV NODE_ENV=production
RUN yarn install --frozen-lockfile --non-interactive --production
CMD ["yarn", "start"]