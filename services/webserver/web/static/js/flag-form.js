$("#flag-form").submit(function (e) {
  e.preventDefault();
  e.stopPropagation();

  $('#correct').removeClass("is-active");
  $('#incorrect').removeClass("is-active");

  $('.circle-loader').toggleClass('is-active');

  data = $(this).serialize();
  data += "&machine-name=" + $('.form-title').text()

  $.ajax({
    type: "POST",
    url: "/machines",
    data: data,
    success: function (response) {
      setTimeout(function () {
        if (response == "success") {
          $('.circle-loader').removeClass('is-active');
          $('#correct').addClass('is-active');
        } else {
          $('.circle-loader').removeClass('is-active');
          $('#incorrect').addClass('is-active');
        }
      }, (.5 * 1000));
    }
  });
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
