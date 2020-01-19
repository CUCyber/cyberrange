$("#machine-creation-form").submit(function (e) {
  e.preventDefault();
  e.stopPropagation();

  form = $(this).FormToJSON();
  form["Points"] = Number(form["Points"]);

  var ws = new WebSocket("ws://" + document.location.host + "/admin/create/machine");

  ws.onopen = function (event) {
    ws.send(JSON.stringify(form));
    spawnProgressBar();
  };

  ws.onerror = function (event) {
    toastr.error('An error occurred. Please reload the page.', 'Machine Form');
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
      toastr.error('Action Request Failed.', 'Machine Form');
      return;
    }

    updateProgressBar(json.Percent);
    switch (json.Type) {
    case "info":
        toastr.info(json.Data, 'Machine Creation Form');
        break;
    case "success":
        toastr.success(json.Data, 'Machine Creation Form');
        break;
    case "error":
        toastr.error(json.Data, 'Machine Creation Form');
        break;
    }
  };
})

$("#machine-deletion-form").submit(function (e) {
  e.preventDefault();
  e.stopPropagation();

  form = $(this).FormToJSON();

  var ws = new WebSocket("ws://" + document.location.host + "/admin/delete/machine");

  ws.onopen = function (event) {
    ws.send(JSON.stringify(form));
    spawnProgressBar();
  };

  ws.onerror = function (event) {
    toastr.error('An error occurred. Please reload the page.', 'Machine Form');
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
      toastr.error('Action Request Failed.', 'Machine Form');
      return;
    }

    updateProgressBar(json.Percent);
    switch (json.Type) {
    case "info":
        toastr.info(json.Data, 'Machine Deletion Form');
        break;
    case "success":
        toastr.success(json.Data, 'Machine Deletion Form');
        break;
    case "error":
        toastr.error(json.Data, 'Machine Deletion Form');
        break;
    }
  };
})
