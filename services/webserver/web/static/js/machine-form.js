$("#machine-form").submit(function (e) {
  e.preventDefault();
  e.stopPropagation();

  data = $(this).serialize();

  spawnNotification("success", "Machine Creation Request Sent.");

  $.ajax({
    type: "POST",
    url: "/admin",
    data: data,
    success: function (response) {
      setTimeout(function () {
        if (response == "success") {
          spawnNotification("success", "Machine Created Successfully");
        } else {
          spawnNotification("danger", response);
        }
      }, (.5 * 1000));
    }
  });
})
