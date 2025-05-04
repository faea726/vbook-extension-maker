load("config.js");

function execute(url, page) {
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
