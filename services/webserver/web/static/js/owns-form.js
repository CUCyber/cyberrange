$("#userown-creation-form").submit(function (e) {
  e.preventDefault();
  e.stopPropagation();

  form = $(this).FormToJSON();

  var ws = new WebSocket("ws://" + document.location.host + "/admin/owns/create_user");

  ws.onopen = function (event) {
    ws.send(JSON.stringify(form));
    spawnProgressBar();
  };

  ws.onerror = function (event) {
    toastr.error('An error occurred. Please reload the page.', 'User Own Creation Form');
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
      toastr.error('Action Request Failed.', 'User Own Creation Form');
      return;
    }

    updateProgressBar(json.Percent);
    switch (json.Type) {
    case "info":
        toastr.info(json.Data, 'User Own Creation Form');
        break;
    case "success":
        toastr.success(json.Data, 'User Own Creation Form');
        break;
    case "error":
        toastr.error(json.Data, 'User Own Creation Form');
        break;
    }
  };
})

$("#rootown-creation-form").submit(function (e) {
  e.preventDefault();
  e.stopPropagation();

  form = $(this).FormToJSON();

  var ws = new WebSocket("ws://" + document.location.host + "/admin/owns/create_root");

  ws.onopen = function (event) {
    ws.send(JSON.stringify(form));
    spawnProgressBar();
  };

  ws.onerror = function (event) {
    toastr.error('An error occurred. Please reload the page.', 'Root Own Creation Form');
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
      toastr.error('Action Request Failed.', 'Root Own Creation Form');
      return;
    }

    updateProgressBar(json.Percent);
    switch (json.Type) {
    case "info":
        toastr.info(json.Data, 'Root Own Creation Form');
        break;
    case "success":
        toastr.success(json.Data, 'Root Own Creation Form');
        break;
    case "error":
        toastr.error(json.Data, 'Root Own Creation Form');
        break;
    }
  };
})

$("#userown-deletion-form").submit(function (e) {
  e.preventDefault();
  e.stopPropagation();

  form = $(this).FormToJSON();

  var ws = new WebSocket("ws://" + document.location.host + "/admin/owns/delete_user");

  ws.onopen = function (event) {
    ws.send(JSON.stringify(form));
    spawnProgressBar();
  };

  ws.onerror = function (event) {
    toastr.error('An error occurred. Please reload the page.', 'User Own Deletion Form');
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
      toastr.error('Action Request Failed.', 'User Own Deletion Form');
      return;
    }

    updateProgressBar(json.Percent);
    switch (json.Type) {
    case "info":
        toastr.info(json.Data, 'User Own Deletion Form');
        break;
    case "success":
        toastr.success(json.Data, 'User Own Deletion Form');
        break;
    case "error":
        toastr.error(json.Data, 'User Own Deletion Form');
        break;
    }
  };
})

$("#rootown-deletion-form").submit(function (e) {
  e.preventDefault();
  e.stopPropagation();

  form = $(this).FormToJSON();

  var ws = new WebSocket("ws://" + document.location.host + "/admin/owns/delete_root");

  ws.onopen = function (event) {
    ws.send(JSON.stringify(form));
    spawnProgressBar();
  };

  ws.onerror = function (event) {
    toastr.error('An error occurred. Please reload the page.', 'Root Own Deletion Form');
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
      toastr.error('Action Request Failed.', 'Root Own Deletion Form');
      return;
    }

    updateProgressBar(json.Percent);
    switch (json.Type) {
    case "info":
        toastr.info(json.Data, 'Root Own Deletion Form');
        break;
    case "success":
        toastr.success(json.Data, 'Root Own Deletion Form');
        break;
    case "error":
        toastr.error(json.Data, 'Root Own Deletion Form');
        break;
    }
  };
})
