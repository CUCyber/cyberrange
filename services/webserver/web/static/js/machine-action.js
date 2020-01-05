let actions = $('#machines-table').find('tbody tr td a');

function actionRequest(action, action_type) {
  $.ajax({
    type: "POST",
    url: "/" + action_type,
    data: "machine-name=" + action.dataset.name,
    success: function (response) {
      setTimeout(function () {
        if (response == "success") {
          spawnNotification("success", "Machine " + action_type + " successful.");
        } else {
          spawnNotification("danger", response);
        }
      }, (.5 * 1000));
    }
  });
}

Array.prototype.forEach.call(actions, function (action) {
  switch (action.dataset.val) {
    case "start":
    case "stop":
    case "revert":
      action.addEventListener('click', function () {
        actionRequest(action, action.dataset.val);
      });
      break;
    default:
      break;
  }
});
