# Finest project out there

## Web app

This is a simple http node server with EJS templates.

### How do I run it locally ?

```bash
cd webapp

# actually runs deno run --allow-net --allow-read --watch main.ts
# this allows the process to access network and reading fils such as the html template
deno task dev
```

And you're done ! The application is running on port 8080.

## Pushing images to docker

```bash
cd webapp

# build it
docker build -t ghcr.io/do3-2023/mmo-kube/webapp:<tag>

# test it
docker run -itp 8080:8080 ghcr.io/do3-2023/mmo-kube/webapp:<tag>

# push it
docker push ghcr.io/do3-2023/mmo-kube/webapp:<tag>
```
