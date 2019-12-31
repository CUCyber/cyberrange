$("#machine-form").submit(function (e) {
  e.preventDefault();
  e.stopPropagation();

  data = $(this).serialize();

  $.ajax({
    type: "POST",
    url: "/admin",
    data: data,
    success: function (response) {
      setTimeout(function () {
        console.log(response);

        var status = document.querySelector('#machine-form-response');
        if (response == "success") {
          status.textContent = "Machine Created Successfully.";
        } else {
          status.textContent = response;
        }

        var modal = document.querySelector('.machine-submission');
        modal.classList.add('is-active');

        ['.modal-background'].forEach(function (e) {
            modal.querySelector(e).addEventListener('click', function (e) {
                e.preventDefault();
                modal.classList.remove('is-active');
            });
        });
      }, (.5 * 1000));
    }
  });
})
