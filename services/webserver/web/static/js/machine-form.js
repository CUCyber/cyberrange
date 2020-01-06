$("#machine-form").submit(function (e) {
  e.preventDefault();
  e.stopPropagation();

  data = $(this).serialize();

  spawnNotification("success", "Machine Creation Request Sent.");

  $.ajax({
    type: "POST",
    url: "/admin",
    data: data,
    timeout: 10 * 60 * 1000,
    success: function (response) {
      setTimeout(function () {
        if (response == "success") {
          spawnNotification("success", "Machine Deployment Started.");
        } else {
          spawnNotification("danger", response);
        }
      }, (.5 * 1000));
    },
  });
})
