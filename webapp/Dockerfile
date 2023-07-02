FROM denoland/deno:alpine-1.34.3 AS build

WORKDIR /app

COPY deps.ts .

RUN deno cache deps.ts

ADD . .

RUN deno cache main.ts

RUN deno bundle main.ts > webapp.bundle.ts

FROM denoland/deno:alpine-1.34.3 AS run

WORKDIR /app

EXPOSE 8080

# Non root
USER deno

COPY --from=build --chown=deno:deno /app/index.ejs .
COPY --from=build --chown=deno:deno /app/webapp.bundle.ts .

CMD [ "run", "--allow-net", "--allow-read", "--allow-env", "webapp.bundle.ts" ]
