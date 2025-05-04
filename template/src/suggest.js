load("config.js");

function execute() {
  return Response.success([
    {
      name: "name",
      link: "link",
      cover: "cover",
      description: "description",
      host: "host",
    },
  ]);
}
