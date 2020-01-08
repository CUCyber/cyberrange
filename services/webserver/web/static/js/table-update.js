var interval;

function copyMachineData(machine, newMachine) {
  machine.ipaddr.innerText = newMachine.IpAddress;
  machine.userowns.innerText = newMachine.UserOwns;
  machine.rootowns.innerText = newMachine.RootOwns;

  $('#machines').find("div.card").each(function(){
    statusObj = $(this).find(".fa-circle")
    statusObj.removeClass(function(index, className) {
        return (className.match("machine-status-.*") || []).join(' ');
    });
    statusObj.addClass("machine-status-" + newMachine.Status);
  });
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

    json.forEach(function (newMachine) {
      $('#machines').find("div.card").each(function(){
        var machine = new Object();

        machine["name"] = $(this).find(".card-header-title")[0];
        machine["ipaddr"] = $(this).find(".machine-address")[0];
        machine["userowns"] = $(this).find(".machine-user-owns")[0];
        machine["rootowns"] = $(this).find(".machine-root-owns")[0];

        if (machine.name.innerText == newMachine.Name) {
          copyMachineData(machine, newMachine);
        }
      });
    });
  });
}

if ($('#machines').length) {
  interval = setInterval(updateMachineList, 5 * 1000);
}
