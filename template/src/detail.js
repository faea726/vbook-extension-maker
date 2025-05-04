load("config.js");

function execute(url) {
  return Response.success({
    name: "name",
    cover: "cover",
    host: "host",
    author: "author",
    description: "description",
    detail: "detail",
    ongoing: true,
    genres: [{ title: "title", input: "input", script: "genrecontent.js" }],
    suggests: [{ title: "title", input: "input", script: "suggest.js" }],
    comments: [{ title: "title", input: "input", script: "comment.js" }],
  });
}
