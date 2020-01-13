function killProgressBar() {
  if ($('#progress-bar').length) {
    $('#progress-bar').remove()
  } 
}

function updateProgressBar(percentage) {
  var progressBar = document.getElementById('progress-bar');

  var curr = progressBar.value;
  var update = setInterval(function() {
    if (curr > percentage) {
      clearInterval(update);
    }
    progressBar.value = curr++;
  }, 15)
}

function spawnProgressBar() {
  if ($('#progress-bar').length) {
    return;
  }

  var progress = document.createElement("progress");

  progress.id = "progress-bar";
  progress.min = 0;
  progress.max = 100;
  progress.classList.add("progress");

  $(".is-12").append(progress);

  interval = setInterval(function () {
    if (progress.value == 100 && progress.parentNode) {
      progress.parentNode.removeChild(progress);
      clearInterval(interval);
    }
  }, 1000);
}
