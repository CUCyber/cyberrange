let machines_table = $('#machines-table');

function copyMachineData(row, machine) {
    row.cells["machine-ip"].innerText = machine.IpAddress;
    row.cells["machine-userowns"].innerText = machine.UserOwns;
    row.cells["machine-rootowns"].innerText = machine.RootOwns;
}

function updateMachineList() {
    $.get('/list', function(resp) {
        let json = JSON.parse(resp);
        let rows = $('#machines-table').find('tbody tr');

        json.forEach(function(machine){
            Array.prototype.forEach.call(rows, function (row){
                if (machine.Name == row.cells["machine-name"].innerText) {
                    copyMachineData(row, machine);
                }
            });
        });
    });
}

if (machines_table.length > 0) {
    setInterval(updateMachineList, 5 * 1000);
}
