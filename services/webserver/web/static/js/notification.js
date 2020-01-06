function spawnNotification(type, body) {
  var div = document.createElement("div");
  var btn = document.createElement("button");

  div.id = "notification";
  div.classList.add("notification");
  div.classList.add("is-" + type);
  div.classList.add("show");
  div.innerText = body;

  btn.classList.add("delete");

  div.appendChild(btn);

  document.body.insertBefore(
    div, document.getElementById("notif-holder")
  );

  ['.delete'].forEach(function (e) {
    div.querySelector(e).addEventListener('click', function (e) {
      div.classList.remove("show");
      div.parentNode.removeChild(div);
    });
  });

  setTimeout(function () {
    if (div.parentNode) {
      div.classList.remove("show");
      div.parentNode.removeChild(div);
    }
  }, 3000);

  console.log(body);
}
