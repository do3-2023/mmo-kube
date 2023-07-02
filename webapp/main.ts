import { config } from "https://deno.land/x/dotenv@v3.2.2/mod.ts";
import { Application, Router, Context } from "./deps.ts";
import { renderFileToString } from "./deps.ts";

interface ApiResponse {
  count: number;
  last_call: Date | null;
  timpestamp: Date;
}

const port = 8080;

// load .env
config({ export: true });

const apiUrl = Deno.env.get("API_URL");

if (!apiUrl) {
  console.error("Error: API_URL environment variable missing");
  Deno.exit(1);
}

const app = new Application();
const router = new Router();

router.get("/", async (ctx: Context) => {
  const response = await fetch(apiUrl);

  const data: ApiResponse = await response.json();

  const myTemplate = await renderFileToString("index.ejs", {
    ...data,
  });

  ctx.response.body = myTemplate;
  ctx.response.status = 200;
  ctx.response.headers.set("Content-Type", "text/html");
});

router.get("/healthz", async (ctx: Context) => {
  try {
    const response = await fetch(`${apiUrl}/healthz`);
    ctx.response.status = response.status;
  } catch (_) {
    ctx.response.status = 500;
  }
});

app.use(router.routes());
app.use(router.allowedMethods());

console.log(`Server is running on http://localhost:${port}`);

await app.listen({ port });
