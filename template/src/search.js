load("config.js");

function execute(key, page) {
  return Response.success(
    [
      {
        name: "name",
        link: "link",
        host: "host",
        cover: "cover",
        description: "description",
      },
    ],
    next
  );
}
