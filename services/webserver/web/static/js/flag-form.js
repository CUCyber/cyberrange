$("#flag-form").submit(function (e) {
  e.preventDefault();
  e.stopPropagation();

  $('#correct').removeClass("is-active");
  $('#incorrect').removeClass("is-active");

  $('.circle-loader').toggleClass('is-active');

  form = $(this).FormToJSON();
  form["Name"] = $('.form-title').text()

  var ws = new WebSocket("ws://" + document.location.host + "/machines/flag");

  ws.onopen = function (event) {
    ws.send(JSON.stringify(form));
  };

  ws.onerror = function (event) {
    toastr.error('An error occurred. Please reload the page.', 'Flag Form');
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
      toastr.error('Action request failed.', 'Flag Form');
      return;
    }

    switch (json.Type) {
      case "info":
        break;
      case "success":
        setTimeout(function () {
          $('.circle-loader').removeClass('is-active');
          $('#correct').addClass('is-active');
        }, (.5 * 1000));
        break;
      case "error":
        setTimeout(function () {
          $('.circle-loader').removeClass('is-active');
          $('#incorrect').addClass('is-active');
        }, (.5 * 1000));
        break;
    }
  };
})

function flagModal(machineName) {
  /* Clear any old flags or load/status icons */
  $("#flag").val('');
  $('#correct').removeClass('is-active');
  $('#incorrect').removeClass('is-active');
  $('.circle-loader').removeClass('is-active');

  var title = document.querySelector('.form-title');
  title.textContent = machineName;

  var modal = document.querySelector('.flag-submission');
  modal.classList.add('is-active');

  ['.modal-background'].forEach(function (e) {
    modal.querySelector(e).addEventListener('click', function (e) {
      e.preventDefault();
      modal.classList.remove('is-active');
    });
  });
}
