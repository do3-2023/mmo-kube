import { serve } from "./deps.ts";
import { renderFileToString } from "./deps.ts";

const port = 8080;

const handler = async (): Promise<Response> => {
  const myTemplate = await renderFileToString("index.ejs", {
    message: "Hello world",
  });
  return new Response(myTemplate, {
    status: 200,
    headers: {
      "Content-Type": "text/html",
    },
  });
};

console.log(`HTTP webserver running. Access it at: http://localhost:8080/`);
await serve(handler, { port });
