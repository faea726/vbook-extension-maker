load("config.js");

function execute() {
  return Response.success([
    { title: "title", input: "input", script: "genrecontent.js" },
  ]);
}
