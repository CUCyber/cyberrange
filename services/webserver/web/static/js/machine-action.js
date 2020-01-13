let actions = $('#machines').find("div.card footer a")

function actionRequest(action, action_type) {
  var ws = new WebSocket("ws://" + Server + "/" + action_type);

  ws.onopen = function (event) {
    ws.send(JSON.stringify({"Name": action.dataset.name}));
    spawnProgressBar();
  };

  ws.onerror = function (event) {
    toastr.error('An error occurred. Please reload the page.', 'Machine Action');
    killProgressBar();
  };

  ws.onmessage = function (event) {
    const json = (function (resp) {
      try {
        return JSON.parse(resp);
      } catch (err) {
        return false;
      }
    })(event.data);

    if (!json) {
      toastr.error('Action request failed.', 'Machine Action');
      return;
    }

    updateProgressBar(json.Percent);
    switch (json.Type) {
    case "info":
        toastr.info(json.Data, 'Machine Action');
        break;
    case "success":
        toastr.success(json.Data, 'Machine Action');
        break;
    case "error":
        toastr.error(json.Data, 'Machine Action');
        break;
    }
  };
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
