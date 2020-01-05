var interval;

function copyMachineData(machine, newMachine) {
  machine.ip.innerText = newMachine.IpAddress;
}

function updateMachineList() {
  $.get('/list', function (resp) {
    const json = (function (resp) {
      try {
        return JSON.parse(resp);
      } catch (err) {
        return false;
      }
    })(resp);

    if (!json) {
      clearInterval(interval);
      return;
    }

    let rows = $('#machines-table').find('tbody tr');

    json.forEach(function (newMachine) {
      Array.prototype.forEach.call(rows, function (row) {
        var machine = new Object();

        Object.keys(row.cells).forEach(function (key) {
          var cell = row.cells[key];
          machine[cell.dataset.val] = cell;
        });

        if (machine.name.innerText == newMachine.Name) {
          copyMachineData(machine, newMachine);
        }
      });
    });
  });
}

if ($('#machines-table').length) {
  interval = setInterval(updateMachineList, 5 * 1000);
}
