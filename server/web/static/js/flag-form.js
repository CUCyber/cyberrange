(function () {
  window.onload = function () {
    var form = document.querySelector("form")
    form.addEventListener("submit", function (e) {
      e.preventDefault()
      var x = new XMLHttpRequest()
      x.onreadystatechange = function () {
        if (x.readyState == 4) {
          console.log(x.response);

          var html = document.querySelector('html');
          html.classList.add('is-clipped');

          if (x.response == "true") {
            var modal = document.querySelector('.correct');
            modal.classList.add('is-active');
          } else {
            var modal = document.querySelector('.incorrect');
            modal.classList.add('is-active');
          }

          ['.delete', '.modal-background', '.modal-close'].forEach(function (e) {
            modal.querySelector(e).addEventListener('click', function (event) {
              event.preventDefault();
              modal.classList.remove('is-active');
              html.classList.remove('is-clipped');
            });
          });
        }
      }
      x.open("POST", window.location.pathname)
      x.send(new FormData(form))
    })
  }
})()
